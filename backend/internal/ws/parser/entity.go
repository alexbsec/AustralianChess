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

type PromotePawnMessage struct {
	Type           string           `json:"type"`
	RequesteeColor chess.PieceColor `json:"requestee_color"`
	PromoteTo      chess.PieceKind  `json:"promote_to"`
	PawnPosition   chess.Position   `json:"pawn_position"`
	DestPosition   chess.Position   `json:"destination_position"`
}
