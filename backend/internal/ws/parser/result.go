package parser

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Result interface {
	IsResult()
}

type FailedCommandResult struct {
	GameStarted bool `json:"game_started"`
}

type MoveResult struct {
	RoomId    string          `json:"room_id"`
	Moved     bool            `json:"moved"`
	GameState chess.GameState `json:"game_state"`
}

type PlayerJoinedResult struct {
	RoomId      string `json:"room_id"`
	PlayerId    string `json:"player_id"`
	Success     bool   `json:"success"`
	GameStarted bool   `json:"game_started"`
}

func (MoveResult) IsResult()          {}
func (PlayerJoinedResult) IsResult()  {}
func (FailedCommandResult) IsResult() {}
