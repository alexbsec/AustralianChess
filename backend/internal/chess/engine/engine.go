package engine

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Engine interface {
	GameResult(gameState *chess.GameState, turnColor chess.PieceColor) *chess.GameResult
	ValidateAndMove(gameState *chess.GameState, requesteeColor chess.PieceColor, fromPos chess.Position, toPos chess.Position) ValidationResponse
	PromotePawn(gameState *chess.GameState, requesteeColor chess.PieceColor, toPiece chess.PieceKind, pawnPos chess.Position, destPos chess.Position) ValidationResponse
}
