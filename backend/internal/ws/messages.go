package ws

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type GameStateMessage struct {
	Type        string            `json:"type"`
	GameState   chess.GameState   `json:"game_state"`
	PlayerColor *chess.PieceColor `json:"player_color,omitempty"`
}

type MovePieceResponse struct {
	Type      string          `json:"type"`
	Moved     bool            `json:"moved"`
	GameState chess.GameState `json:"game_state"`
}

type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
