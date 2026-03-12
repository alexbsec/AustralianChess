package ws

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type GameStateMessage struct {
	Type        string            `json:"type"`
	GameState   chess.GameState   `json:"gameState"`
	PlayerColor *chess.PieceColor `json:"playerColor,omitempty"`
}

type MovePieceResponse struct {
	Type      string          `json:"type"`
	Moved     bool            `json:"moved"`
	GameState chess.GameState `json:"gameState"`
}

type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
