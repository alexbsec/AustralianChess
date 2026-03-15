package engine

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Arbiter interface {
	CanMovePiece(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string)
	IsInCheck(board chess.Board, turnColor chess.PieceColor) bool
	IsCheckmate(board chess.Board, turnColor chess.PieceColor) bool
	IsStalemate(board chess.Board, turnColor chess.PieceColor) bool
}
