package parser

import (
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

type Command interface {
	IsCommand()
}

type MoveCommand struct {
	RoomId         string
	RequesteeColor chess.PieceColor
	FromPos        chess.Position
	ToPos          chess.Position
}

func (MoveCommand) IsCommand() {}
