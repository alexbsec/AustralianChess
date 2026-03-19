package engine

import (
	"fmt"
	"log"

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

func (ce *ChessEngine) GameResult(gameState *chess.GameState, turnColor chess.PieceColor) *chess.GameResult {
	if gameState == nil {
		return nil
	}

	if ce.arbiter.IsCheckmate(gameState.Board, turnColor) {
		result := chess.ResultCheckmate
		return &result
	}

	if ce.arbiter.IsStalemate(gameState.Board, turnColor) {
		result := chess.ResultStalemate
		return &result
	}

	return nil
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

	opponentColor := chess.OponentColor(requesteeColor)
	if ce.arbiter.IsInCheck(gameState.Board, opponentColor) {
		gameState.InCheck = true
	} else {
		gameState.InCheck = false
	}

	switchTurn(gameState)
	return ValidationResponse{
		Type:    SuccessResponse,
		Message: feedback,
	}
}

func (ce *ChessEngine) PromotePawn(gameState *chess.GameState, requesteeColor chess.PieceColor, toPiece chess.PieceKind, pawnPos chess.Position, destPos chess.Position) ValidationResponse {
	if gameState == nil {
		return ValidationResponse{
			Type:    ErrorResponse,
			Message: "no engine was passed to validate the move",
		}
	}

	piece := gameState.Board.GetPiece(pawnPos)
	if piece == nil {
		return ValidationResponse{
			Type:    FailureResponse,
			Message: fmt.Sprintf("no piece to move at position: %v", pawnPos),
		}
	}

	if piece.Kind != chess.PawnPiece {
		return ValidationResponse{
			Type: FailureResponse,
			Message: fmt.Sprintf("cannot promote non-pawn piece. Want: %d, got: %d", chess.PawnPiece, piece.Kind),
		}
	}

	if piece.Color != requesteeColor {
		return ValidationResponse{
			Type: FailureResponse,
			Message: "cannot promote opponent's pawn piece",
		}
	}


	log.Printf("pawnPos: %v, destPos: %v", pawnPos, destPos)
	canMove, feedback := ce.arbiter.CanMovePiece(gameState.Board, *piece, pawnPos, destPos)
	if !canMove {
		return ValidationResponse{
			Type: FailureResponse,
			Message: fmt.Sprintf("cannot move pawn to promotion site. From: %v, To: %v", pawnPos, destPos),
		}
	}

	err := movePiece(gameState, pawnPos, destPos)
	if err != nil {
		return ValidationResponse{
			Type:    ErrorResponse,
			Message: "cannot move on empty square",
		}
	}

	opponentColor := chess.OponentColor(requesteeColor)
	if ce.arbiter.IsInCheck(gameState.Board, opponentColor) {
		gameState.InCheck = true
	} else {
		gameState.InCheck = false
	}

	promotePawn(gameState, requesteeColor, toPiece, destPos)	
	switchTurn(gameState)

	return ValidationResponse{
		Type:    SuccessResponse,
		Message: feedback,
	}
}

func promotePawn(gameState *chess.GameState, requesteeColor chess.PieceColor, toPiece chess.PieceKind, pawnPos chess.Position) {
	promotedPiece := chess.NewPiece(toPiece, requesteeColor)
	nextBoard := gameState.Board.Clone()
	nextBoard.Data[pawnPos.Row][pawnPos.Col].Piece = promotedPiece
	gameState.Board = nextBoard
}

func movePiece(gameState *chess.GameState, fromPos, toPos chess.Position) error {
	nextBoard := gameState.Board.Clone()
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
