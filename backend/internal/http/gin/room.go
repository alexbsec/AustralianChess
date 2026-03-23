package gin

import (
	"context"
	"net/http"

	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/gin-gonic/gin"
)

func RoomHandler(
	ctx context.Context,
	routerGroup *gin.RouterGroup,
	roomService rooms.IService,
	wsHandler *ws.Handler,
) {
	routerGroup.Handle("GET", "/create", MakeNewRoom(roomService))
	routerGroup.Handle("GET", "/:id", ServeRoom(roomService))
	if wsHandler != nil {
		routerGroup.Handle("GET", "/ws/:id", wsHandler.HandleRoom)
	}
}

func MakeNewRoom(roomService rooms.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.MustGet("sessionId").(int64)
		_ = ctx.MustGet("userId").(int64)

		res, err := roomService.CreateRoom(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"roomId": res.RoomId})
	}
}

func ServeRoom(roomService rooms.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.MustGet("sessionId").(int64)
		_ = ctx.MustGet("userId").(int64)

		id := ctx.Param("id")
		res, err := roomService.DisplayRoom(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}
