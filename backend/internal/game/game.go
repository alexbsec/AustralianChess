package game

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	wsTypes "github.com/alexbsec/AustralianChess/backend/internal/ws/ws_types"
)

type Game interface {
	HandlePlayBotRoom(ctx context.Context, roomId string, playerColor chess.PieceColor, difficulty bot.Difficulty) error
	HandleConnect(ctx context.Context, roomId string, client *wsTypes.Client) (contracts.Result, chess.GameState, error)
	HandleCommand(ctx context.Context, roomId string, client *wsTypes.Client, cmd contracts.Command) (contracts.Result, error)
	HandleDisconnect(roomId string, client *wsTypes.Client)
	HandleDestroyRoom(ctx context.Context, roomId string) error
}
