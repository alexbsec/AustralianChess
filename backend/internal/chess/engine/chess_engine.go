package engine

import (
	"fmt"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

type ChessEngine struct {
	arbiter Arbiter
}

func NewChessEngine() *ChessEngine {
	return &ChessEngine{
		arbiter: NewChessArbitrator(),
	}
}

func (ce *ChessEngine) ValidateAndMove(gameState *chess.GameState, requesteeColor chess.PieceColor, fromPos chess.Position, toPos chess.Position) ValidationResponse {
	if gameState == nil {
		return ValidationResponse{
			Type:    ErrorResponse,
			Message: "no engine was passed to validate the move",
		}
	}

	piece := gameState.Board.GetPiece(fromPos)
	if piece == nil {
		return ValidationResponse{
			Type:    FailureResponse,
			Message: fmt.Sprintf("no piece to move at position: %v", fromPos),
		}
	}

	turn := gameState.Turn
	if turn != requesteeColor {
		return ValidationResponse{
			Type:    FailureResponse,
			Message: "cannot move on oponent's turn",
		}
	}

	if piece.Color != requesteeColor {
		return ValidationResponse{
			Type:    FailureResponse,
			Message: "cannot move oponent's piece",
		}
	}

	canMove, feedback := ce.arbiter.CanMovePiece(gameState.Board, *piece, fromPos, toPos)
	if !canMove {
		return ValidationResponse{
			Type:    FailureResponse,
			Message: feedback,
		}
	}

	err := movePiece(gameState, fromPos, toPos)
	if err != nil {
		return ValidationResponse{
			Type:    ErrorResponse,
			Message: "cannot move on empty square",
		}
	}

	switchTurn(gameState)
	return ValidationResponse{
		Type:    SuccessResponse,
		Message: feedback,
	}
}

func movePiece(gameState *chess.GameState, fromPos, toPos chess.Position) error {
	nextBoard := gameState.Board
	err := nextBoard.MovePiece(fromPos, toPos)
	if err != nil {
		return err
	}

	gameState.Board = nextBoard
	return nil
}

func switchTurn(gameState *chess.GameState) {
	if gameState.Turn == chess.PieceWhite {
		gameState.Turn = chess.PieceBlack
		return
	}

	gameState.Turn = chess.PieceWhite
}
