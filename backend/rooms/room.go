package rooms

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
)

type IService interface {
	CreateRoom(ctx context.Context) (CreateRoomResponse, error)
	DisplayRoom(ctx context.Context, roomId string) (DisplayRoomResponse, error)
	FetchRoom(ctx context.Context, roomId string) (*Room, error)
	DeleteRoom(ctx context.Context, roomId string) error

	ExecuteCommand(ctx context.Context, cmd contracts.Command) (contracts.Result, error)
	UpdatePlayerJoined(ctx context.Context, roomId, playerId string, color chess.PieceColor) (contracts.Result, error)
}
