package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/alexbsec/AustralianChess/backend/internal/ws/parser"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const bufferSize = 1024

type Handler struct {
	roomService rooms.IService
	hub         IHub
	upgrader    websocket.Upgrader
	wsParser    parser.IParser
}

func NewHandler(roomService rooms.IService, hub IHub) *Handler {
	wsParser := parser.NewRoomParser()

	return &Handler{
		roomService: roomService,
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
}

func (h *Handler) HandleRoom(ctx *gin.Context) {
	roomId := ctx.Param("id")
	if roomId == "" {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{
			Type:    "error",
			Message: "missing room id",
		})
		return
	}

	playerId := ctx.Query("playerId")
	if playerId == "" {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{
			Type:    "error",
			Message: "missing player",
		})
	}

	reqCtx := ctx.Request.Context()
	room, err := h.roomService.FetchRoom(reqCtx, roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{
			Type:    "error",
			Message: err.Error(),
		})
		return
	}

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}

	client, err := h.hub.AddClient(roomId, playerId, conn)
	if err != nil {
		log.Printf("failed to add client to hub: %v", err)
		_ = conn.Close()
		return
	}

	defer func() {
		h.hub.RemoveClient(roomId, conn)
		_ = conn.Close()
	}()

	if err := conn.WriteJSON(GameStateMessage{
		Type:        "game_state",
		GameState:   *room.GameState,
		PlayerColor: client.Color,
	}); err != nil {
		log.Printf("failed to write initial game state: %v", err)
		return
	}

	h.loop(reqCtx, roomId, client)
}

func (h *Handler) loop(ctx context.Context, roomId string, client *Client) {
	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message: %v", err)
			return
		}

		if client.Role == Spectator {
			if err := client.Conn.WriteJSON(ErrorMessage{
				Type:    "error",
				Message: "spectators cannot send commands",
			}); err != nil {
				log.Printf("failed to write spectator error: %v", err)
				return
			}
			continue
		}

		cmd, err := h.wsParser.ParseMessage(roomId, data)
		if err != nil {
			log.Printf("failed to parse message: %v", err)
			if err := client.Conn.WriteJSON(ErrorMessage{
				Type:    "error",
				Message: err.Error(),
			}); err != nil {
				log.Printf("failed to write parse error: %v", err)
				return
			}
			continue
		}

		cmd = cmd.SetPlayerId(client.PlayerId)
		cmd = cmd.SetRequesterColor(*client.Color)

		result, err := h.roomService.ExecuteCommand(ctx, cmd)
		if err != nil {
			log.Printf("failed to execute command: %v", err)
			if err := client.Conn.WriteJSON(ErrorMessage{
				Type:    "error",
				Message: err.Error(),
			}); err != nil {
				log.Printf("failed to write execute error: %v", err)
				return
			}
			continue
		}

		if err := h.hub.Broadcast(roomId, result); err != nil {
			log.Printf("failed to broadcast response: %v", err)
			return
		}
	}
}
