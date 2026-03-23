package engine_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/engine"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/engine/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newTestGameState() *chess.GameState {
	board := chess.NewInitialBoard()

	return &chess.GameState{
		Board: board,
		Turn:  chess.PieceWhite,
	}
}

func TestGameResult_Checkmate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)

	game := newTestGameState()

	mockArbiter.EXPECT().
		IsCheckmate(game.Board, chess.PieceWhite).
		Return(true)

	eng := engine.NewChessEngineWithArbiter(mockArbiter)
	result := eng.GameResult(game, chess.PieceWhite)

	require.NotNil(t, result)
	require.Equal(t, chess.ResultCheckmate, *result)
}


func TestGameResult_Stalemate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)
	game := newTestGameState()

	mockArbiter.EXPECT().
		IsCheckmate(game.Board, chess.PieceWhite).
		Return(false)

	mockArbiter.EXPECT().
		IsStalemate(game.Board, chess.PieceWhite).
		Return(true)

	eng := engine.NewChessEngineWithArbiter(mockArbiter)

	result := eng.GameResult(game, chess.PieceWhite)

	require.NotNil(t, result)
	require.Equal(t, chess.ResultStalemate, *result)
}


func TestValidateAndMove_NoGameState(t *testing.T) {
	eng := engine.NewChessEngine()

	resp := eng.ValidateAndMove(nil, chess.PieceWhite, chess.Position{}, chess.Position{})

	require.Equal(t, engine.ErrorResponse, resp.Type)
	require.Contains(t, resp.Message, "no engine was passed")
}


func TestValidateAndMove_NoPiece(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)
	eng := engine.NewChessEngineWithArbiter(mockArbiter) 
	game := newTestGameState()

	resp := eng.ValidateAndMove(game, chess.PieceWhite,
		chess.Position{Row: 5, Col: 5},
		chess.Position{Row: 6, Col: 5},
	)

	require.Equal(t, engine.FailureResponse, resp.Type)
	require.Contains(t, resp.Message, "no piece to move")
}


func TestValidateAndMove_WrongTurn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)
	eng := engine.NewChessEngineWithArbiter(mockArbiter)

	game := newTestGameState()
	game.Turn = chess.PieceBlack

	resp := eng.ValidateAndMove(game, chess.PieceWhite,
		chess.Position{Row: 1, Col: 0},
		chess.Position{Row: 2, Col: 0},
	)

	require.Equal(t, engine.FailureResponse, resp.Type)
	require.Contains(t, resp.Message, "oponent")
}


func TestValidateAndMove_InvalidMove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)
	eng := engine.NewChessEngineWithArbiter(mockArbiter)

	game := newTestGameState()

	mockArbiter.EXPECT().
		CanMovePiece(game.Board, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, "illegal move")

	resp := eng.ValidateAndMove(game, chess.PieceWhite,
		chess.Position{Row: 1, Col: 0},
		chess.Position{Row: 2, Col: 0},
	)

	require.Equal(t, engine.FailureResponse, resp.Type)
	require.Equal(t, "illegal move", resp.Message)
}


func TestValidateAndMove_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)
	eng := engine.NewChessEngineWithArbiter(mockArbiter)

	game := newTestGameState()

	// Ensure piece exists
	game.Board.Data[1][0].Piece = chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	mockArbiter.EXPECT().
		CanMovePiece(game.Board, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, "ok")

	mockArbiter.EXPECT().
		IsInCheck(gomock.Any(), chess.PieceBlack).
		Return(false)

	resp := eng.ValidateAndMove(game, chess.PieceWhite,
		chess.Position{Row: 1, Col: 0},
		chess.Position{Row: 2, Col: 0},
	)

	require.Equal(t, engine.SuccessResponse, resp.Type)
	require.Equal(t, chess.PieceBlack, game.Turn)
}


func TestPromotePawn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArbiter := mocks.NewMockArbiter(ctrl)
	eng := engine.NewChessEngineWithArbiter(mockArbiter)

	game := newTestGameState()

	// Place pawn manually at promotion source
	game.Board.Data[1][0].Piece = chess.NewPiece(chess.PawnPiece, chess.PieceWhite)

	mockArbiter.EXPECT().
		CanMovePiece(game.Board, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, "promoted")

	mockArbiter.EXPECT().
		IsInCheck(gomock.Any(), chess.PieceBlack).
		Return(false)

	resp := eng.PromotePawn(
		game,
		chess.PieceWhite,
		chess.QueenPiece,
		chess.Position{Row: 1, Col: 0},
		chess.Position{Row: 0, Col: 0},
	)

	require.Equal(t, engine.SuccessResponse, resp.Type)

	// Promotion should result in new piece
	p := game.Board.GetPiece(chess.Position{Row: 0, Col: 0})
	require.NotNil(t, p)
	require.Equal(t, chess.QueenPiece, p.Kind)
	require.Equal(t, chess.PieceWhite, p.Color)
}
