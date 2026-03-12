package parser

import (
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

type Command interface {
	SetPlayerId(playerId string) Command
	SetRequesterColor(color chess.PieceColor) Command
	GetRoomId() string
}

type MoveCommand struct {
	RoomId         string
	PlayerId       string
	RequesteeColor chess.PieceColor
	FromPos        chess.Position
	ToPos          chess.Position
}

func (m MoveCommand) SetPlayerId(playerId string) Command {
	m.PlayerId = playerId
	return m
}

func (m MoveCommand) SetRequesterColor(color chess.PieceColor) Command {
	m.RequesteeColor = color
	return m
}

func (m MoveCommand) GetRoomId() string {
	return m.RoomId
}
