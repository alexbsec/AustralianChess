package rooms

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/internal/ws/parser"
)

type IService interface {
	CreateRoom(ctx context.Context) (CreateRoomResponse, error)
	DisplayRoom(ctx context.Context, roomId string) (DisplayRoomResponse, error)
	FetchRoom(ctx context.Context, roomId string) (*Room, error)

	ExecuteCommand(ctx context.Context, cmd parser.Command) (parser.Result, error)
}
