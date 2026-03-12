package parser

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Result interface {
	IsResult()
}

type MoveResult struct {
	RoomId    string
	Moved     bool
	GameState chess.GameState
}

func (MoveResult) IsResult() {}
