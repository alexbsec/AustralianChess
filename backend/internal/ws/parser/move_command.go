package parser

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

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

func (m MoveCommand) SetColor(color chess.PieceColor) Command {
	m.RequesteeColor = color
	return m
}

func (m MoveCommand) GetRoomId() string {
	return m.RoomId
}
