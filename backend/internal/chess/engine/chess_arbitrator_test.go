package engine_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/engine"
	"github.com/stretchr/testify/require"
)

const boardSize = 12

func emptyCleanBoard() chess.Board {
	board := chess.Board{
		Data: make([][]chess.Square, boardSize),
	}

	for r := range boardSize {
		board.Data[r] = make([]chess.Square, boardSize)
	}

	whiteKing := chess.NewPiece(chess.KingPiece, chess.PieceWhite)
	blackKing := chess.NewPiece(chess.KingPiece, chess.PieceBlack)
	placePiece(&board, whiteKing, chess.Position{Row: 0, Col: 0})
	placePiece(&board, blackKing, chess.Position{Row: 11, Col: 11})

	return board
}

func placePiece(b *chess.Board, piece *chess.Piece, pos chess.Position) {
	b.Data[pos.Row][pos.Col].Piece = piece
}

func TestCanMovePiece_OutOfBounds(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	piece := chess.NewPiece(chess.KnightPiece, chess.PieceWhite)
	placePiece(&board, piece, chess.Position{Row: 0, Col: 0})

	ok, msg := arb.CanMovePiece(
		board,
		*piece,
		chess.Position{Row: 0, Col: 0},
		chess.Position{Row: -1, Col: 0},
	)

	require.False(t, ok)
	require.Contains(t, msg, "movement not within board limits")
}

func TestCanMovePiece_SameSquare(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()
	piece := chess.NewPiece(chess.KnightPiece, chess.PieceWhite)

	placePiece(&board, piece, chess.Position{Row: 0, Col: 0})

	ok, msg := arb.CanMovePiece(
		board,
		*piece,
		chess.Position{Row: 0, Col: 0},
		chess.Position{Row: 0, Col: 0},
	)

	require.False(t, ok)
	require.Contains(t, msg, "distance is 0")
}

func TestCanMovePiece_KnightValid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()
	piece := chess.NewPiece(chess.KnightPiece, chess.PieceWhite)

	placePiece(&board, piece, chess.Position{Row: 4, Col: 4})

	ok, _ := arb.CanMovePiece(
		board,
		*piece,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 6, Col: 5},
	)

	require.True(t, ok)
}

func TestCanMovePiece_KnightInvalid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()
	piece := chess.NewPiece(chess.KnightPiece, chess.PieceWhite)

	placePiece(&board, piece, chess.Position{Row: 4, Col: 4})

	ok, msg := arb.CanMovePiece(
		board,
		*piece,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 5, Col: 5},
	)

	require.False(t, ok)
	require.Contains(t, msg, "cannot move")
}

func TestCanMovePiece_QueenValidMoves(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	queen := chess.NewPiece(chess.QueenPiece, chess.PieceWhite)
	placePiece(&board, queen, chess.Position{Row: 4, Col: 4})

	tests := []struct {
		name string
		to   chess.Position
	}{
		{"horizontal", chess.Position{Row: 4, Col: 7}},
		{"vertical", chess.Position{Row: 7, Col: 4}},
		{"diagonal", chess.Position{Row: 7, Col: 7}},
		{"diagonal_reverse", chess.Position{Row: 1, Col: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := arb.CanMovePiece(board, *queen,
				chess.Position{Row: 4, Col: 4},
				tt.to,
			)

			require.True(t, ok)
		})
	}
}

func TestCanMovePiece_RookBlocked(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)
	blocker := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	placePiece(&board, rook, chess.Position{Row: 0, Col: 0})
	placePiece(&board, blocker, chess.Position{Row: 0, Col: 2})

	ok, _ := arb.CanMovePiece(
		board,
		*rook,
		chess.Position{Row: 0, Col: 0},
		chess.Position{Row: 0, Col: 5},
	)

	require.False(t, ok)
}

func TestCanMovePiece_RookValid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	rook := chess.NewPiece(chess.RookPiece, chess.PieceWhite)

	placePiece(&board, rook, chess.Position{Row: 1, Col: 0})

	ok, _ := arb.CanMovePiece(
		board,
		*rook,
		chess.Position{Row: 1, Col: 0},
		chess.Position{Row: 1, Col: 5},
	)

	require.True(t, ok)
}

func TestCanMovePiece_BishopValid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	bishop := chess.NewPiece(chess.BishopPiece, chess.PieceWhite)

	placePiece(&board, bishop, chess.Position{Row: 4, Col: 4})

	ok, _ := arb.CanMovePiece(
		board,
		*bishop,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 6, Col: 6},
	)

	require.True(t, ok)
}

func TestCanMovePiece_KingValid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	king := chess.NewPiece(chess.KingPiece, chess.PieceWhite)

	placePiece(&board, king, chess.Position{Row: 4, Col: 4})

	ok, _ := arb.CanMovePiece(
		board,
		*king,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 5, Col: 5},
	)

	require.True(t, ok)
}

func TestCanMovePiece_DiagonalCapture_Valid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	enemy := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)

	from := chess.Position{Row: 10, Col: 1}
	to := chess.Position{Row: 9, Col: 2}

	placePiece(&board, pawn, from)
	placePiece(&board, enemy, to)

	ok, _ := arb.CanMovePiece(board, *pawn, from, to)

	require.True(t, ok)
}

func TestCanMovePiece_PawnForward(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	placePiece(&board, pawn, chess.Position{Row: 1, Col: 1})

	ok, _ := arb.CanMovePiece(
		board,
		*pawn,
		chess.Position{Row: 1, Col: 1},
		chess.Position{Row: 0, Col: 1},
	)

	require.True(t, ok)
}

func TestCanMovePiece_TwoSquare_NotFromStart(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	// NOT starting row anymore
	from := chess.Position{Row: 2, Col: 4}
	to := chess.Position{Row: 4, Col: 4}

	placePiece(&board, pawn, from)

	ok, _ := arb.CanMovePiece(board, *pawn, from, to)

	require.False(t, ok)
}


func TestCanMovePiece_TwoSquareMove_Valid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	// starting position
	from := chess.Position{Row: 10, Col: 4}
	to := chess.Position{Row: 8, Col: 4}

	placePiece(&board, pawn, from)

	ok, _ := arb.CanMovePiece(board, *pawn, from, to)

	require.True(t, ok)
}


func TestCanMovePiece_TwoSquareBlocked(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	pawn := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)
	blocker := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	from := chess.Position{Row: 1, Col: 4}
	to := chess.Position{Row: 3, Col: 4}

	placePiece(&board, pawn, from)
	placePiece(&board, blocker, chess.Position{Row: 2, Col: 4}) // blocking square

	ok, _ := arb.CanMovePiece(board, *pawn, from, to)

	require.False(t, ok)
}

func TestCanMovePiece_KangarooValid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	kangaroo := chess.NewPiece(chess.KangarooPiece, chess.PieceWhite)
	placePiece(&board, kangaroo, chess.Position{Row: 4, Col: 4})

	// 2 horizontal
	ok, _ := arb.CanMovePiece(board, *kangaroo,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 4, Col: 6},
	)
	require.True(t, ok)

	// 3 vertical
	ok, _ = arb.CanMovePiece(board, *kangaroo,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 7, Col: 4},
	)
	require.True(t, ok)
}

func TestCanMovePiece_KangarooInvalid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	kangaroo := chess.NewPiece(chess.KangarooPiece, chess.PieceWhite)
	placePiece(&board, kangaroo, chess.Position{Row: 4, Col: 4})

	ok, msg := arb.CanMovePiece(board, *kangaroo,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 5, Col: 5},
	)

	require.False(t, ok)
	require.Contains(t, msg, "cannot move")
}

func TestCanMovePiece_KangarooCapture(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	kangaroo := chess.NewPiece(chess.KangarooPiece, chess.PieceWhite)
	enemy := chess.NewPiece(chess.PawnPiece, chess.PieceBlack)

	placePiece(&board, kangaroo, chess.Position{Row: 4, Col: 4})
	placePiece(&board, enemy, chess.Position{Row: 4, Col: 6})

	ok, _ := arb.CanMovePiece(board, *kangaroo,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 4, Col: 6},
	)

	require.True(t, ok)
}

func TestCanMovePiece_OligarchAsKing(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	oligarch := chess.NewPiece(chess.OligarchPiece, chess.PieceWhite)
	placePiece(&board, oligarch, chess.Position{Row: 4, Col: 4})

	ok, _ := arb.CanMovePiece(board, *oligarch,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 5, Col: 5},
	)

	require.True(t, ok)
}

func TestCanMovePiece_OligarchAsQueen(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	oligarch := chess.NewPiece(chess.OligarchPiece, chess.PieceWhite)
	friend := chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	placePiece(&board, oligarch, chess.Position{Row: 4, Col: 4})
	placePiece(&board, friend, chess.Position{Row: 4, Col: 5}) // activates queen mode

	ok, _ := arb.CanMovePiece(board, *oligarch,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 7, Col: 4}, // rook-like
	)

	require.True(t, ok)
}

func TestCanMovePiece_OligarchInvalid(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	oligarch := chess.NewPiece(chess.OligarchPiece, chess.PieceWhite)
	placePiece(&board, oligarch, chess.Position{Row: 4, Col: 4})

	ok, msg := arb.CanMovePiece(board, *oligarch,
		chess.Position{Row: 4, Col: 4},
		chess.Position{Row: 6, Col: 7}, // invalid for both modes
	)

	require.False(t, ok)
	require.Contains(t, msg, "cannot move")
}

func TestIsInCheck(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	king := chess.NewPiece(chess.KingPiece, chess.PieceWhite)
	rook := chess.NewPiece(chess.RookPiece, chess.PieceBlack)

	placePiece(&board, king, chess.Position{Row: 4, Col: 4})
	placePiece(&board, rook, chess.Position{Row: 4, Col: 0})

	require.True(t, arb.IsInCheck(board, chess.PieceWhite))
}

func TestIsCheckmate(t *testing.T) {
	arb := engine.NewChessArbitrator()

	board := emptyCleanBoard()

	rook1 := chess.NewPiece(chess.RookPiece, chess.PieceBlack)
	rook2 := chess.NewPiece(chess.RookPiece, chess.PieceBlack)

	// Attack king directly
	placePiece(&board, rook1, chess.Position{Row: 0, Col: 10}) // horizontal check
	placePiece(&board, rook2, chess.Position{Row: 1, Col: 11}) // vertical check

	require.True(t, arb.IsCheckmate(board, chess.PieceWhite))
}

