package gin

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/gin-gonic/gin"
)

func MakeHandlers(
	ctx context.Context,
	roomService rooms.IService,
	userService users.IService,
	wsHandler *ws.Handler,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	routerGroup := router.Group("/api/v1")
	routerGroup.Use(gin.Logger())

	if roomService != nil {
		RoomHandler(ctx, routerGroup, roomService, wsHandler)
	}

	if userService != nil {
		UserHandler(ctx, routerGroup, userService)
	}

	return router
}
