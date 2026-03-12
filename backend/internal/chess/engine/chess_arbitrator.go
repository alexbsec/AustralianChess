package engine

import (
	"fmt"
	"log"
	"slices"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

var successMoveResFmt string = "piece %d movement accepted from %v to %v"
var failureMoveResFmt string = "piece %d cannot move from %v to %v"
var captureMoveResFmt string = "piece %d capture piece %d on movement from %v to %v"

type ChessArbitrator struct{}

func NewChessArbitrator() *ChessArbitrator {
	return &ChessArbitrator{}
}

func (ca *ChessArbitrator) CanMovePiece(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	if !isInsideBoard(board, fromPos, toPos) {
		return false, "movement not within board limits"
	}

	if isSameSquare(fromPos, toPos) {
		return false, "movement distance is 0"
	}

	var canMove bool
	var feedback string

	switch piece.Kind {
	case chess.PawnPiece:
		canMove, feedback = canMovePawn(board, piece, fromPos, toPos)
	case chess.RookPiece:
		canMove, feedback = canMoveRook(board, piece, fromPos, toPos)
	case chess.KnightPiece:
		canMove, feedback = canMoveKnight(board, piece, fromPos, toPos)
	case chess.BishopPiece:
		canMove, feedback = canMoveBishop(board, piece, fromPos, toPos)
	case chess.KangarooPiece:
		canMove, feedback = canMoveKangaroo(board, piece, fromPos, toPos)
	case chess.QueenPiece:
		canMove, feedback = canMoveQueen(board, piece, fromPos, toPos)
	case chess.KingPiece:
		canMove, feedback = canMoveKing(board, piece, fromPos, toPos)
	case chess.OligarchPiece:
		canMove, feedback = canMoveOligarch(board, piece, fromPos, toPos)
	default:
		return false, "unknown piece"
	}

	if !canMove {
		return false, feedback
	}

	nextBoard, err := applyMoveToBoard(board, fromPos, toPos)
	if err != nil {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	kingInCheck, checkFeedback := isKingInCheck(nextBoard, piece.Color)
	if kingInCheck {
		return false, checkFeedback
	}

	return true, feedback
}

func canMoveOligarch(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	hasFriendlyNeighbor := hasAdjacentFriendlyPiece(board, piece.Color, fromPos)
	if hasFriendlyNeighbor {
		return canMoveQueen(board, piece, fromPos, toPos)
	}

	return canMoveKing(board, piece, fromPos, toPos)
}

func canMoveKing(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	rowDiff := intAbs(toPos.Row - fromPos.Row)
	colDiff := intAbs(toPos.Col - fromPos.Col)

	if rowDiff > 1 || colDiff > 1 {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	targetPiece := board.GetPiece(toPos)
	if targetPiece == nil {
		return true, successMovementMessage(piece.Kind, fromPos, toPos)
	}

	if targetPiece.Color == piece.Color {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	return true, successCaptureMessage(piece.Kind, targetPiece.Kind, fromPos, toPos)
}

func canMoveQueen(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	if ok, msg := canMoveRook(board, piece, fromPos, toPos); ok {
		return true, msg
	}

	if ok, msg := canMoveBishop(board, piece, fromPos, toPos); ok {
		return true, msg
	}

	return false, failureMovementMessage(piece.Kind, fromPos, toPos)
}

func canMoveKangaroo(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	rowDiff := intAbs(toPos.Row - fromPos.Row)
	colDiff := intAbs(toPos.Col - fromPos.Col)

	isKangarooMove := (rowDiff == 2 && colDiff == 0) ||
		(rowDiff == 0 && colDiff == 2) ||
		(rowDiff == 3 && colDiff == 0) ||
		(rowDiff == 0 && colDiff == 3)
	if !isKangarooMove {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	targetPiece := board.GetPiece(toPos)
	if targetPiece == nil {
		return true, successMovementMessage(piece.Kind, fromPos, toPos)
	}

	if targetPiece.Color == piece.Color {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	return true, successCaptureMessage(piece.Kind, targetPiece.Kind, fromPos, toPos)
}

func canMoveBishop(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	rowDiff := intAbs(toPos.Row - fromPos.Row)
	colDiff := intAbs(toPos.Col - fromPos.Col)

	if rowDiff == 0 || rowDiff != colDiff {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	rowStep := -1
	if toPos.Row > fromPos.Row {
		rowStep = 1
	}

	colStep := -1
	if toPos.Col > fromPos.Col {
		colStep = 1
	}

	row := fromPos.Row + rowStep
	col := fromPos.Col + colStep

	for row != toPos.Row {
		pos := chess.Position{Row: row, Col: col}
		if board.GetPiece(pos) != nil {
			return false, failureMovementMessage(piece.Kind, fromPos, toPos)
		}

		row += rowStep
		col += colStep
	}

	targetPiece := board.GetPiece(toPos)
	if targetPiece == nil {
		return true, successMovementMessage(piece.Kind, fromPos, toPos)
	}

	if targetPiece.Color == piece.Color {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	return true, successCaptureMessage(piece.Kind, targetPiece.Kind, fromPos, toPos)
}

func canMoveKnight(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	rowDiff := intAbs(toPos.Row - fromPos.Row)
	colDiff := intAbs(toPos.Col - fromPos.Col)

	isKnightMove := (rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2)
	if !isKnightMove {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	targetPiece := board.GetPiece(toPos)
	if targetPiece == nil {
		return true, successMovementMessage(piece.Kind, fromPos, toPos)
	}

	if targetPiece.Color == piece.Color {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	return true, successCaptureMessage(piece.Kind, targetPiece.Kind, fromPos, toPos)
}

func canMoveRook(
	board chess.Board,
	piece chess.Piece,
	fromPos, toPos chess.Position,
) (bool, string) {
	isSameRow := fromPos.Row == toPos.Row
	isSameCol := fromPos.Col == toPos.Col

	if !isSameRow && !isSameCol {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	rowStep := 0
	if toPos.Row > fromPos.Row {
		rowStep = 1
	} else if toPos.Row < fromPos.Row {
		rowStep = -1
	}

	colStep := 0
	if toPos.Col > fromPos.Col {
		colStep = 1
	} else if toPos.Col < fromPos.Col {
		colStep = -1
	}

	row := fromPos.Row + rowStep
	col := fromPos.Col + colStep

	for row != toPos.Row || col != toPos.Col {
		if board.GetPiece(chess.Position{Row: row, Col: col}) != nil {
			return false, failureMovementMessage(piece.Kind, fromPos, toPos)
		}

		row += rowStep
		col += colStep
	}

	targetPiece := board.GetPiece(toPos)
	if targetPiece == nil {
		return true, successMovementMessage(piece.Kind, fromPos, toPos)
	}

	if targetPiece.Color == piece.Color {
		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	return true, successCaptureMessage(piece.Kind, targetPiece.Kind, fromPos, toPos)
}

func canMovePawn(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	direction := -1
	if piece.Color == chess.PieceBlack {
		direction = 1
	}

	rowDiff := toPos.Row - fromPos.Row
	colDiff := toPos.Col - fromPos.Col
	targetPiece := board.GetPiece(toPos)

	// Forward movement
	if fromPos.Col == toPos.Col {
		if rowDiff == direction && targetPiece == nil {
			return true, successMovementMessage(piece.Kind, fromPos, toPos)
		}

		// condition to move 2 squares at starting position
		condition := isPieceAtStartingPosition(board, piece.Kind, piece.Color, fromPos) &&
			rowDiff == 2*direction && targetPiece == nil && board.Data[fromPos.Row+direction][fromPos.Col].Piece == nil
		if condition {
			return true, successMovementMessage(piece.Kind, fromPos, toPos)
		}

		return false, failureMovementMessage(piece.Kind, fromPos, toPos)
	}

	// diagonal capture
	if intAbs(colDiff) == 1 && rowDiff == direction && targetPiece != nil && targetPiece.Color != piece.Color {
		return true, successCaptureMessage(piece.Kind, targetPiece.Kind, fromPos, toPos)
	}

	return false, failureMovementMessage(piece.Kind, fromPos, toPos)
}

func isInsideBoard(board chess.Board, fromPos, toPos chess.Position) bool {
	rows, cols := board.Dimensions()
	if fromPos.Row < 0 || fromPos.Row >= rows || fromPos.Col < 0 || fromPos.Col >= cols ||
		toPos.Row < 0 || toPos.Row >= rows || toPos.Col < 0 || toPos.Col >= cols {
		return false
	}

	return true
}

func isSameSquare(fromPos, toPos chess.Position) bool {
	if fromPos.Col == toPos.Col && fromPos.Row == toPos.Row {
		return true
	}

	return false
}

func isPieceAtStartingPosition(board chess.Board, pieceKind chess.PieceKind, pieceColor chess.PieceColor, position chess.Position) bool {
	rows, cols := board.Dimensions()
	backRow := rows - 1
	if pieceColor == chess.PieceBlack {
		backRow = 0
	}

	if pieceKind == chess.PawnPiece {
		if backRow == 0 {
			return position.Row == backRow+1
		}
		return position.Row == backRow-1
	}

	// not pawn, so must be at back row
	if position.Row != backRow {
		return false
	}

	startingColumnsByPiece := map[chess.PieceKind][]int{
		chess.RookPiece:     {0, cols - 1},
		chess.OligarchPiece: {1, cols - 2},
		chess.KnightPiece:   {2, cols - 3},
		chess.KangarooPiece: {3, cols - 4},
		chess.BishopPiece:   {4, cols - 5},
		chess.QueenPiece:    {5},
		chess.KingPiece:     {6},
	}

	possibleColumns, ok := startingColumnsByPiece[pieceKind]
	if !ok {
		// this should never hit
		log.Printf("piece not found on possible pieces: %d", pieceKind)
		return false
	}

	return slices.Contains(possibleColumns, position.Col)
}

func isKingInCheck(board chess.Board, pieceColor chess.PieceColor) (bool, string) {
	kingPosition := findKingPosition(board, pieceColor)
	if kingPosition == nil {
		return true, fmt.Sprintf("king not found for color: %d", pieceColor)
	}

	oponentColor := chess.OponentColor(pieceColor)
	kingSquareAttacked := isSquareAttacked(board, *kingPosition, oponentColor)
	if kingSquareAttacked {
		return true, "cannot move piece that will leave king in check"
	}

	return false, ""
}

func isSquareAttacked(board chess.Board, squarePosition chess.Position, color chess.PieceColor) bool {
	rows, cols := board.Dimensions()
	for row := range rows {
		for col := range cols {
			piecePos := chess.Position{Row: row, Col: col}
			piece := board.GetPiece(piecePos)
			if piece == nil || piece.Color != color {
				continue
			}

			if canPieceAttackSquare(board, *piece, piecePos, squarePosition) {
				return true
			}
		}
	}

	return false
}

func canPieceAttackSquare(board chess.Board, piece chess.Piece, piecePos, squarePos chess.Position) bool {
	if !isInsideBoard(board, piecePos, squarePos) {
		return false
	}

	if isSameSquare(piecePos, squarePos) {
		return false
	}

	canMove := false
	switch piece.Kind {
	case chess.PawnPiece:
		direction := -1
		if piece.Color == chess.PieceBlack {
			direction = 1
		}
		rowDiff := squarePos.Row - piecePos.Row
		colDiff := squarePos.Col - piecePos.Col
		canMove = rowDiff == direction && intAbs(colDiff) == 1
	case chess.RookPiece:
		canMove, _ = canMoveRook(board, piece, piecePos, squarePos)
	case chess.KnightPiece:
		canMove, _ = canMoveKnight(board, piece, piecePos, squarePos)
	case chess.BishopPiece:
		canMove, _ = canMoveBishop(board, piece, piecePos, squarePos)
	case chess.KangarooPiece:
		canMove, _ = canMoveKangaroo(board, piece, piecePos, squarePos)
	case chess.OligarchPiece:
		canMove, _ = canMoveOligarch(board, piece, piecePos, squarePos)
	case chess.QueenPiece:
		canMove, _ = canMoveQueen(board, piece, piecePos, squarePos)
	case chess.KingPiece:
		canMove, _ = canMoveKing(board, piece, piecePos, squarePos)
	default:
		break
	}

	return canMove
}

func findKingPosition(board chess.Board, pieceColor chess.PieceColor) *chess.Position {
	rows, cols := board.Dimensions()
	for row := range rows {
		for col := range cols {
			position := &chess.Position{Row: row, Col: col}
			piece := board.GetPiece(*position)
			if piece != nil && piece.Kind == chess.KingPiece && piece.Color == pieceColor {
				return position
			}
		}
	}

	return nil
}

func hasAdjacentFriendlyPiece(board chess.Board, color chess.PieceColor, position chess.Position) bool {
	for rowOffset := -1; rowOffset <= 1; rowOffset++ {
		for colOffset := -1; colOffset <= 1; colOffset++ {
			if rowOffset == 0 && colOffset == 0 {
				continue
			}

			nextRow := position.Row + rowOffset
			nextCol := position.Col + colOffset
			nextPos := chess.Position{Row: nextRow, Col: nextCol}

			if !isInsideBoard(board, position, nextPos) {
				continue
			}

			neighbor := board.GetPiece(nextPos)
			if neighbor != nil && neighbor.Color == color {
				return true
			}
		}
	}

	return false
}

func applyMoveToBoard(board chess.Board, fromPos, toPos chess.Position) (chess.Board, error) {
	nextBoard := board.Clone()
	err := nextBoard.MovePiece(fromPos, toPos)
	return nextBoard, err
}

func successMovementMessage(kind chess.PieceKind, fromPos, toPos chess.Position) string {
	return fmt.Sprintf(successMoveResFmt, kind, fromPos, toPos)
}

func failureMovementMessage(kind chess.PieceKind, fromPos, toPos chess.Position) string {
	return fmt.Sprintf(failureMoveResFmt, kind, fromPos, toPos)
}

func successCaptureMessage(attacker chess.PieceKind, attacked chess.PieceKind, fromPos, toPos chess.Position) string {
	return fmt.Sprintf(captureMoveResFmt, attacker, attacked, fromPos, toPos)
}

func intAbs(val int) int {
	if val < 0 {
		return -1 * val
	}

	return val
}
