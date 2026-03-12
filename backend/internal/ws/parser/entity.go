package parser

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Envelope struct {
	Type string `json:"type"`
}

type MovePieceMessage struct {
	Type           string           `json:"type"`
	RequesteeColor chess.PieceColor `json:"requesteeColor"`
	FromPos        chess.Position   `json:"fromPos"`
	ToPos          chess.Position   `json:"toPos"`
}
