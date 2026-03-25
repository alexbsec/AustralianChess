package chess_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/stretchr/testify/require"
)

func TestPieceColorString_White(t *testing.T) {
	result := chess.PieceColorString(chess.PieceWhite)
	require.Equal(t, "white", result)
}

func TestPieceColorString_Black(t *testing.T) {
	result := chess.PieceColorString(chess.PieceBlack)
	require.Equal(t, "black", result)
}

func TestPieceColorString_Unknown(t *testing.T) {
	result := chess.PieceColorString(chess.PieceColor(99))
	require.Equal(t, "unknown", result)
}
