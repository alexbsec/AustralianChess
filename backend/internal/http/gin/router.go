package gin

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/gin-gonic/gin"
)

func MakeHandlers(
	ctx context.Context,
	roomService rooms.IService,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	routerGroup := router.Group("/api/v1")
	routerGroup.Use(gin.Logger())

	if roomService != nil {
		RoomHandler(ctx, routerGroup, roomService)
	}

	return router
}
