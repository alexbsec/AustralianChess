package gin

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/internal/auth"
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
	authorizer auth.IAuthorizer,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	authRouter := router.Group("/api/v1/room")
	generalRouterGroup := router.Group("/api/v1")
	authRouter.Use(AuthMiddleware(authorizer))
	generalRouterGroup.Use(gin.Logger())

	if roomService != nil {
		RoomHandler(ctx, authRouter, roomService, wsHandler)
	}

	if userService != nil {
		UserHandler(ctx, generalRouterGroup, userService)
	}

	return router
}
