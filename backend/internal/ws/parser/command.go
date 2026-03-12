package parser

import (
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

type Command interface {
	SetPlayerId(playerId string) Command
	SetColor(color chess.PieceColor) Command
	GetRoomId() string
}
