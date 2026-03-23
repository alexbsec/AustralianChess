package contracts

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type PromoteCommand struct {
	RoomId         string
	PlayerId       string
	RequesteeColor chess.PieceColor
	PromoteTo      chess.PieceKind
	PiecePosition  chess.Position
	DestPosition   chess.Position
}

func (pc PromoteCommand) SetPlayerId(playerId string) Command {
	pc.PlayerId = playerId
	return pc
}

func (pc PromoteCommand) SetColor(color chess.PieceColor) Command {
	pc.RequesteeColor = color
	return pc
}

func (pc PromoteCommand) GetRoomId() string {
	return pc.RoomId
}
