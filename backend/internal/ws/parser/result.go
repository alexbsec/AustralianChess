package parser

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Result interface {
	IsResult()
}

type FailedCommandResult struct {
	GameStarted bool `json:"gameStarted"`
}

type MoveResult struct {
	RoomId    string          `json:"roomId"`
	Moved     bool            `json:"moved"`
	GameState chess.GameState `json:"gameState"`
}

type PlayerJoinedResult struct {
	RoomId      string `json:"roomId"`
	PlayerId    string `json:"playerId"`
	Success     bool   `json:"success"`
	GameStarted bool   `json:"gameStarted"`
}

func (MoveResult) IsResult()          {}
func (PlayerJoinedResult) IsResult()  {}
func (FailedCommandResult) IsResult() {}
