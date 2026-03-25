package bot

import (
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	"github.com/alexbsec/AustralianChess/backend/rooms"
)

type Manager interface {
	SpawnBot(roomId string, color chess.PieceColor,	difficulty Difficulty, roomService rooms.IService, onMove func(string, contracts.Result) error) error
	NotifyState(roomId string, state chess.GameState)
	RemoveBot(roomId string)
	HasBot(roomId string) bool
}
