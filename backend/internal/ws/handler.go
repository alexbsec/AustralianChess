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

	h.hub.AddClient(roomId, conn)
	defer func() {
		h.hub.RemoveClient(roomId, conn)
		_ = conn.Close()
	}()

	if err := conn.WriteJSON(GameStateMessage{
		Type:      "game_state",
		GameState: *room.GameState,
	}); err != nil {
		log.Printf("failed to write initial game state: %v", err)
		return
	}

	h.loop(ctx, roomId, conn)
}

func (h *Handler) loop(ctx context.Context, roomId string, conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message: %v", err)
			return
		}

		cmd, err := h.wsParser.ParseMessage(roomId, data)
		if err != nil {
			log.Printf("failed to parse message: %v", err)
			return
		}

		result, err := h.roomService.ExecuteCommand(ctx, cmd)
		if err != nil {
			log.Printf("failed to execute command: %v", err)
			return
		}

		if err := conn.WriteJSON(result); err != nil {
			log.Printf("failed to write response: %v", err)
			return
		}
	}
}
