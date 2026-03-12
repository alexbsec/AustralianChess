package main

import (
	"context"

	ginChess "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	"github.com/alexbsec/AustralianChess/backend/rooms"
)

func main() {
	ctx := context.Background()
	roomSvc := rooms.NewService()
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(roomSvc, hub)
	router := ginChess.MakeHandlers(ctx, roomSvc, wsHandler)
	router.Run(":8080")
}
