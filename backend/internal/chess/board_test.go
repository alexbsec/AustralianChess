package chess_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/stretchr/testify/require"
)

func TestNewInitialBoard(t *testing.T) {
	board := chess.NewInitialBoard()

	rows, cols := board.Dimensions()
	require.Equal(t, 12, rows)
	require.Equal(t, 12, cols)
}


func TestBoardColors(t *testing.T) {
	board := chess.NewInitialBoard()

	for r := range 12 {
		for c := range 12 {
			expected := chess.SquareColorLight
			if (r+c)%2 == 1 {
				expected = chess.SquareColorDark
			}

			require.Equal(t, expected, board.Data[r][c].Color)
		}
	}
}


func TestInitialPiecesPlacement(t *testing.T) {
	board := chess.NewInitialBoard()

	// Black back rank
	require.NotNil(t, board.Data[0][0].Piece)
	require.Equal(t, chess.PieceBlack, board.Data[0][0].Piece.Color)

	// White back rank
	require.NotNil(t, board.Data[11][0].Piece)
	require.Equal(t, chess.PieceWhite, board.Data[11][0].Piece.Color)

	// Black pawns
	require.NotNil(t, board.Data[1][0].Piece)
	require.Equal(t, chess.PieceBlack, board.Data[1][0].Piece.Color)

	// White pawns
	require.NotNil(t, board.Data[10][0].Piece)
	require.Equal(t, chess.PieceWhite, board.Data[10][0].Piece.Color)
}

func TestBoardClone(t *testing.T) {
	board := chess.NewInitialBoard()
	cloned := board.Clone()

	// Modify clone
	cloned.Data[0][0].Piece = nil

	// Original must remain unchanged
	require.NotNil(t, board.Data[0][0].Piece)

	// Ensure different memory references (deep copy)
	require.NotSame(t, &board.Data[0][0], &cloned.Data[0][0])
}

func TestMovePiece_Success(t *testing.T) {
	board := chess.NewInitialBoard()

	from := chess.Position{Row: 1, Col: 0} // pawn
	to := chess.Position{Row: 2, Col: 0}

	err := board.MovePiece(from, to)
	require.NoError(t, err)

	require.Nil(t, board.Data[from.Row][from.Col].Piece)
	require.NotNil(t, board.Data[to.Row][to.Col].Piece)
}

func TestMovePiece_NoPiece(t *testing.T) {
	board := chess.NewInitialBoard()

	from := chess.Position{Row: 5, Col: 5}
	to := chess.Position{Row: 6, Col: 5}

	err := board.MovePiece(from, to)

	require.Error(t, err)
	require.Contains(t, err.Error(), "no piece to move")
}


func TestDimensions(t *testing.T) {
	board := chess.NewInitialBoard()

	rows, cols := board.Dimensions()

	require.Equal(t, 12, rows)
	require.Equal(t, 12, cols)
}
