package parser

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Envelope struct {
	Type string `json:"type"`
}

type MovePieceMessage struct {
	Type           string           `json:"type"`
	RequesteeColor chess.PieceColor `json:"requestee_color"`
	FromPos        chess.Position   `json:"from_pos"`
	ToPos          chess.Position   `json:"to_pos"`
}
