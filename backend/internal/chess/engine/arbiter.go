package engine

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Arbiter interface {
	CanMovePiece(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string)
}
