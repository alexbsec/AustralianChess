package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/alexbsec/AustralianChess/backend/internal/game"
	"github.com/alexbsec/AustralianChess/backend/internal/ws/parser"
	wsTypes "github.com/alexbsec/AustralianChess/backend/internal/ws/ws_types"
	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const bufferSize = 1024

type Handler struct {
	userService      users.IService
	hub              IHub
	upgrader         websocket.Upgrader
	wsParser         parser.IParser
	gameOrchestrator game.Game
}

func NewHandler(userService users.IService, hub IHub, gameOrchestrator game.Game) *Handler {
	wsParser := parser.NewRoomParser()

	h := &Handler{
		userService: userService,
		hub:         hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  bufferSize,
			WriteBufferSize: bufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		wsParser:   wsParser,
		gameOrchestrator: gameOrchestrator,
	}

	h.hub.SetOnRoomEmptyCallback(func(roomId string) {
		err := h.gameOrchestrator.HandleDestroyRoom(context.Background(), roomId)
		if err != nil {
			log.Printf("failed to delete empty room %s: %v", roomId, err)
		}
	})

	h.hub.StartJanitor()
	return h
}

func (h *Handler) PlayBot(ctx *gin.Context, playBotDTO PlayBotDTO) {
	userId := h.authorize(ctx)

	reqCtx := ctx.Request.Context()
	user, err := h.userService.User(reqCtx, userId)
	if err != nil {
		log.Printf("failed to fetch user %d: %v", userId, err)
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{
			Type:    "error",
			Message: "failed to fetch user",
		})
		return
	}

	roomId := playBotDTO.RoomId
	playerColor := playBotDTO.PlayerPlayingAs
	difficulty := playBotDTO.Difficulty

	err = h.gameOrchestrator.HandlePlayBotRoom(reqCtx, roomId, playerColor, difficulty)
	if err != nil {
		log.Printf("failed to create bot room: %v", err)
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{
			Type: "error",
			Message: "failed to create bot room",
		})
		return
	}
	log.Printf("bot added to room %s", roomId)

	botId := bot.DifficultyToBotId(difficulty)
	botColor := chess.OponentColor(playerColor)
	h.hub.AddBot(roomId, botId, botColor)

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{
			Type: "error",
			Message: "failed to upgrade connection",
		})
		return
	}
	log.Printf("bot room connection upgraded")

	client, err := h.hub.AddClient(roomId, user.PlayerId, conn)
	if err != nil {
		log.Printf("failed to add client to hub: %v", err)
		_ = conn.Close()
		return
	}
	log.Printf("client added to hub")

	defer func() {
		h.gameOrchestrator.HandleDisconnect(roomId, client)
		h.hub.RemoveClient(roomId, conn)
		_ = conn.Close()
	}()

	result, initialState, err := h.gameOrchestrator.HandleConnect(reqCtx, roomId, client)
	if err != nil {
		log.Printf("failed to handle connect for room %s: %v", roomId, err)
        _ = client.WriteJSON(ErrorMessage{Type: "error", Message: "connection was not successful"})
        return
	}
	log.Printf("initialState: %v", initialState)

    if err := client.WriteJSON(GameStateMessage{
        Type:        "game_state",
        GameState:   initialState,
        PlayerColor: client.Color,
    }); err != nil {
        log.Printf("failed to write initial game state: %v", err)
        return
    }
	log.Printf("message sent containing initialState")

    if result != nil {
        if err := h.hub.Broadcast(roomId, result, true); err != nil {
            log.Printf("failed to broadcast player joining: %v", err)
            return
        }
    }

	log.Printf("client %v added to hub for room %s and playerId %s with role %s", client, roomId, user.PlayerId, client.Role)
	h.loop(reqCtx, roomId, client)
}

func (h *Handler) HandleRoom(ctx *gin.Context, roomId string) {
	log.Printf("handling new websocket connection for room %s", ctx.Param("id"))
	userId := h.authorize(ctx)

	reqCtx := ctx.Request.Context()
	user, err := h.userService.User(reqCtx, userId)
	if err != nil {
		log.Printf("failed to fetch user %d: %v", userId, err)
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{
			Type:    "error",
			Message: "failed to fetch user",
		})
		return
	}
	playerId := user.PlayerId

	log.Printf("upgrading connection to websocket for room %s and playerId %s", roomId, playerId)
	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}

	log.Printf("adding client to hub for room %s and playerId %s", roomId, playerId)
	client, err := h.hub.AddClient(roomId, playerId, conn)
	if err != nil {
		log.Printf("failed to add client to hub: %v", err)
		_ = conn.Close()
		return
	}

	log.Printf("client %v added to hub for room %s and playerId %s with role %s", client, roomId, playerId, client.Role)
	defer func() {
		h.gameOrchestrator.HandleDisconnect(roomId, client)
		h.hub.RemoveClient(roomId, conn)
		_ = conn.Close()
	}()

	result, initialState, err := h.gameOrchestrator.HandleConnect(reqCtx, roomId, client)
	if err != nil {
		log.Printf("failed to handle connection for room %s: %v", roomId, err)
		_ = client.WriteJSON(ErrorMessage{Type: "error", Message: "connection was not successful"})
		return
	}

	// handle connect
    if err := client.WriteJSON(GameStateMessage{
        Type:        "game_state",
        GameState:   initialState,
        PlayerColor: client.Color,
    }); err != nil {
        log.Printf("failed to write initial game state: %v", err)
        return
    }

	if result != nil {
		if err := h.hub.Broadcast(roomId, result, true); err != nil {
			log.Printf("failed to broadcast player joining: %v", err)
		}
	}

	h.loop(reqCtx, roomId, client)
}

func (h *Handler) loop(ctx context.Context, roomId string, client *wsTypes.Client) {
	go keepAlivePing(ctx, client)
	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			log.Printf("client %s disconnected from room %s", client.PlayerId, roomId)
			return
		}

		if client.Role == wsTypes.Spectator {
			client.WriteJSON(ErrorMessage{Type: "error", Message: "spectators cannot send commands"})
			continue
		}

		cmd, err := h.wsParser.ParseMessage(roomId, data)
		if err != nil {
			log.Printf("parser failed to parse data: %v, err: %v", string(data), err)
			client.WriteJSON(ErrorMessage{Type: "error", Message: "could not parse message"})
			continue
		}

		result, err := h.gameOrchestrator.HandleCommand(ctx, roomId, client, cmd)
		if err != nil {
			log.Printf("game orchestrator failed to handle command. Err: %v", err)
			client.WriteJSON(ErrorMessage{Type: "error", Message: "internal server error"})
			continue
		}

		if err := h.hub.Broadcast(roomId, result, true); err != nil {
			log.Printf("hub failed to broadcast result. Err: %v", err)
			client.WriteJSON(ErrorMessage{Type: "error", Message: "failed to broadcast result"})
		}
	}
}

func (h *Handler) authorize(ctx *gin.Context) int64 {
	_ = ctx.MustGet("sessionId").(int64)
	userId := ctx.MustGet("userId").(int64)
	return userId
}

func keepAlivePing(ctx context.Context, client *wsTypes.Client) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("sending ping to client %s", client.PlayerId)
			if err := client.WritePing(); err != nil {
				log.Printf("ping failed for client %s: %v", client.PlayerId, err)
				return
			}
			log.Printf("ping sent to client %s", client.PlayerId)
		case <-client.Done:
			log.Printf("keepalive stopped for client %s (done)", client.PlayerId)
			return
		case <-ctx.Done():
			log.Printf("keepalive stopped for client %s (ctx)", client.PlayerId)
			return
		}
	}
}
