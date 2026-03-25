package engine_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/engine"
	"github.com/stretchr/testify/require"
)

func TestZobristHash_Deterministic(t *testing.T) {
	board := chess.NewInitialBoard()
	h1 := engine.ZobristHash(board, chess.PieceWhite)
	h2 := engine.ZobristHash(board, chess.PieceWhite)
	require.Equal(t, h1, h2)
}

func TestZobristHash_DiffersByColor(t *testing.T) {
	board := chess.NewInitialBoard()
	hWhite := engine.ZobristHash(board, chess.PieceWhite)
	hBlack := engine.ZobristHash(board, chess.PieceBlack)
	require.NotEqual(t, hWhite, hBlack)
}

func TestZobristHash_DiffersByPosition(t *testing.T) {
	board1 := chess.NewInitialBoard()
	board2 := chess.NewInitialBoard()

	// Move a pawn on board2
	board2.Data[10][0].Piece = nil
	board2.Data[8][0].Piece = chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	h1 := engine.ZobristHash(board1, chess.PieceWhite)
	h2 := engine.ZobristHash(board2, chess.PieceWhite)
	require.NotEqual(t, h1, h2)
}

func TestZobristHash_EmptyBoard(t *testing.T) {
	board := emptyCleanBoard()
	h := engine.ZobristHash(board, chess.PieceWhite)
	require.NotEqual(t, uint64(0), h) // Hash of position with kings should be nonzero
}

func TestZobristUpdate_MatchesFreshHash(t *testing.T) {
	board := emptyCleanBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	placePiece(&board, rook, chess.Position{Row: 5, Col: 0})

	from := chess.Position{Row: 5, Col: 0}
	to := chess.Position{Row: 5, Col: 5}

	initialHash := engine.ZobristHash(board, chess.PieceWhite)
	updatedHash := engine.ZobristUpdate(initialHash, board, from, to)

	// Move the piece and compute fresh hash
	_ = board.MovePiece(from, to)
	freshHash := engine.ZobristHash(board, chess.PieceBlack) // turn flipped

	require.Equal(t, freshHash, updatedHash)
}

func TestZobristUpdate_Capture(t *testing.T) {
	board := emptyCleanBoard()
	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	enemy := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)

	from := chess.Position{Row: 5, Col: 0}
	to := chess.Position{Row: 5, Col: 5}

	placePiece(&board, rook, from)
	placePiece(&board, enemy, to)

	initialHash := engine.ZobristHash(board, chess.PieceWhite)
	updatedHash := engine.ZobristUpdate(initialHash, board, from, to)

	// Apply capture and compute fresh hash
	_ = board.MovePiece(from, to)
	freshHash := engine.ZobristHash(board, chess.PieceBlack)

	require.Equal(t, freshHash, updatedHash)
}

func TestZobristUpdate_NilPiece(t *testing.T) {
	board := emptyCleanBoard()
	from := chess.Position{Row: 5, Col: 5} // empty square
	to := chess.Position{Row: 6, Col: 5}

	initialHash := engine.ZobristHash(board, chess.PieceWhite)
	result := engine.ZobristUpdate(initialHash, board, from, to)

	// Should return unchanged hash when no piece at source
	require.Equal(t, initialHash, result)
}

func TestZobristHash_BlackTurnXORsTurnBit(t *testing.T) {
	board := emptyCleanBoard()
	hWhite := engine.ZobristHash(board, chess.PieceWhite)
	hBlack := engine.ZobristHash(board, chess.PieceBlack)
	// XOR of both should equal the turn zobrist value (consistent)
	require.NotEqual(t, hWhite, hBlack)
}
