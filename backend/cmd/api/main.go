package main

import (
	"context"

	ginChess "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	"github.com/alexbsec/AustralianChess/backend/rooms"
)

func main() {
	ctx := context.Background()

	roomSvc := rooms.NewService()

	router := ginChess.MakeHandlers(ctx, roomSvc)
	router.Run(":8080")
}
