package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/ws/parser"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const bufferSize = 1024

type Handler struct {
	roomService rooms.IService
	userService users.IService
	hub         IHub
	upgrader    websocket.Upgrader
	wsParser    parser.IParser
}

func NewHandler(roomService rooms.IService, userService users.IService, hub IHub) *Handler {
	wsParser := parser.NewRoomParser()

	h := &Handler{
		roomService: roomService,
		userService: userService,
		hub:         hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  bufferSize,
			WriteBufferSize: bufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		wsParser: wsParser,
	}

	h.hub.SetOnRoomEmptyCallback(func(roomId string) {
		roomService := h.roomService
		err := roomService.DeleteRoom(context.Background(), roomId)
		if err != nil {
			log.Printf("failed to delete empty room %s: %v", roomId, err)
		}
	})

	h.hub.StartJanitor()
	return h
}

func (h *Handler) HandleRoom(ctx *gin.Context) {
	log.Printf("handling new websocket connection for room %s", ctx.Param("id"))
	_ = ctx.MustGet("sessionId").(int64)
	userId := ctx.MustGet("userId").(int64)

	roomId := ctx.Param("id")
	if roomId == "" {
		log.Printf("missing room id in request")
		ctx.JSON(http.StatusBadRequest, ErrorMessage{
			Type:    "error",
			Message: "missing room id",
		})
		return
	}

	log.Printf("fetching user %d for websocket connection", userId)

	user, err := h.userService.User(ctx.Request.Context(), userId)
	if err != nil {
		log.Printf("failed to fetch user %d: %v", userId, err)
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{
			Type:    "error",
			Message: "failed to fetch user",
		})
		return
	}
	playerId := user.PlayerId

	log.Printf("fetching room %s for websocket connection", roomId)
	reqCtx := ctx.Request.Context()
	room, err := h.roomService.FetchRoom(reqCtx, roomId)
	if err != nil {
		log.Printf("REJECTED: Room %s not found in Service: %v", roomId, err) // CHECK THIS LOG
		ctx.JSON(http.StatusBadRequest, ErrorMessage{
			Type:    "error",
			Message: err.Error(),
		})
		return
	}

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
		h.hub.RemoveClient(roomId, conn)
		_ = conn.Close()
	}()

	log.Printf("client connected to room %s with playerId %s", roomId, playerId)
	if err := client.WriteJSON(GameStateMessage{
		Type:        "game_state",
		GameState:   *room.GameState,
		PlayerColor: client.Color,
	}); err != nil {
		log.Printf("failed to write initial game state: %v", err)
		return
	}

	log.Printf("client role in room %s is %s", roomId, client.Role)
	if client.Role == PlayerOne || client.Role == PlayerTwo {
		log.Printf("updating room %s with player %s joining as %s", roomId, playerId, client.Role)
		result, err := h.roomService.UpdatePlayerJoined(ctx, roomId, playerId, *client.Color)
		if err != nil {
			log.Printf("failed to player join game: %v", err)
			return
		}

		log.Printf("broadcasting player joining room %s as %s", roomId, client.Role)
		if err := h.hub.Broadcast(roomId, result, true); err != nil {
			log.Printf("failed to broadcast player joining room: %v", err)
			return
		}
	}

	log.Printf("client connected to room %s as %s", roomId, client.Role)
	h.loop(reqCtx, roomId, client)
}

func (h *Handler) loop(ctx context.Context, roomId string, client *Client) {
	go keepAlivePong(ctx, client)
    for {
        _, data, err := client.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            log.Printf("failed to read message: %v, data: %v", err, data)
            return
        }

        if client.Role == Spectator {
            if err := client.WriteJSON(ErrorMessage{Type: "error", Message: "spectators cannot send commands"}); err != nil {
                log.Printf("failed to write spectator error: %v", err)
                return
            }
            continue
        }

        cmd, err := h.wsParser.ParseMessage(roomId, data)
        if err != nil {
            log.Printf("failed to parse message: %v", err)
            if err := client.WriteJSON(ErrorMessage{Type: "error", Message: err.Error()}); err != nil {
                log.Printf("failed to write parse error: %v", err)
                return
            }
            continue
        }

        cmd = cmd.SetPlayerId(client.PlayerId)
        cmd = cmd.SetColor(*client.Color)

        result, err := h.roomService.ExecuteCommand(ctx, cmd)
        if err != nil {
            log.Printf("failed to execute command: %v", err)
            if err := client.WriteJSON(ErrorMessage{Type: "error", Message: err.Error()}); err != nil {
                log.Printf("failed to write execute error: %v", err)
                return
            }
            continue
        }

        if err := h.hub.Broadcast(roomId, result, true); err != nil {
            log.Printf("failed to broadcast response: %v", err)
            return
        }
    }
}

func keepAlivePong(ctx context.Context, client *Client) {
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
