package chess_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	board := chess.NewInitialBoard()
	game := chess.NewGame(board)

	require.NotNil(t, game)
	require.Equal(t, chess.PieceWhite, game.Turn)
	require.Nil(t, game.Result)
	require.Nil(t, game.EndReason)
	require.False(t, game.InCheck)
	require.Nil(t, game.PendingPromotion)
}

func TestOponentColor_WhiteReturnsBlack(t *testing.T) {
	result := chess.OponentColor(chess.PieceWhite)
	require.Equal(t, chess.PieceBlack, result)
}

func TestOponentColor_BlackReturnsWhite(t *testing.T) {
	result := chess.OponentColor(chess.PieceBlack)
	require.Equal(t, chess.PieceWhite, result)
}
