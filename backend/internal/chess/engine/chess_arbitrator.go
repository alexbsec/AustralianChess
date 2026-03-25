package engine

import (
	"fmt"
	"log"
	"slices"
	"sync/atomic"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)


var successMoveResFmt string = "piece %d movement accepted from %v to %v"
var failureMoveResFmt string = "piece %d cannot move from %v to %v"
var captureMoveResFmt string = "piece %d capture piece %d on movement from %v to %v"

type ChessArbitrator struct {
	blackPiecesPositions []chess.PiecePosition
	whitePiecesPositions []chess.PiecePosition
	tt                   map[uint64]ttEntry
	stop                 atomic.Bool
}

func NewChessArbitrator() *ChessArbitrator {
	return &ChessArbitrator{
		tt: make(map[uint64]ttEntry, 1<<16),
	}
}

// SetStop allows external callers to signal the search to stop early.
func (ca *ChessArbitrator) SetStop(v bool) {
	ca.stop.Store(v)
}

// IsStopped reports whether a stop was requested.
func (ca *ChessArbitrator) IsStopped() bool {
	return ca.stop.Load()
}

func (ca *ChessArbitrator) CanMovePiece(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	if !chess.IsInsideBoard(board, fromPos, toPos) {
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

func (ca *ChessArbitrator) IsInCheck(board chess.Board, turnColor chess.PieceColor) bool {
	isCheck, _ := isKingInCheck(board, turnColor)
	return isCheck
}

func (ca *ChessArbitrator) IsCheckmate(board chess.Board, turnColor chess.PieceColor) bool {
	legalMoves := ca.getLegalMovesForColor(board, turnColor)
	isKingInCheck, _ := isKingInCheck(board, turnColor)

	if len(legalMoves) == 0 && isKingInCheck {
		return true
	}

	return false
}

func (ca *ChessArbitrator) IsStalemate(board chess.Board, turnColor chess.PieceColor) bool {
	legalMoves := ca.getLegalMovesForColor(board, turnColor)
	isKingInCheck, _ := isKingInCheck(board, turnColor)

	if len(legalMoves) == 0 && !isKingInCheck {
		return true
	}

	return false
}

func (ca *ChessArbitrator) BestMove(board chess.Board, color chess.PieceColor, depth int) (chess.Move, bool) {
    legalMoves := ca.getLegalMovesForColor(board, color)
    if len(legalMoves) == 0 {
        return chess.Move{}, false
    }
    legalMoves = ca.orderMoves(board, legalMoves)

    hash := ZobristHash(board, color)

    bestMove := legalMoves[0]
    alpha, beta := minScore, maxScore
    bestScore := minScore

    for _, move := range legalMoves {
        if ca.stop.Load() {
            return chess.Move{}, false
        }

        next, err := applyMoveToBoard(board, move.FromPos, move.ToPos)
        if err != nil {
            continue
        }

        nextHash := ZobristUpdate(hash, board, move.FromPos, move.ToPos)
        score := -ca.alphaBeta(next, depth-1, -beta, -alpha, chess.OponentColor(color), nextHash)

        if score > bestScore {
            bestScore = score
            bestMove = move
        }
        if score > alpha {
            alpha = score
        }
    }

    if ca.stop.Load() {
        return chess.Move{}, false
    }
    return bestMove, true
}

func (ca *ChessArbitrator) alphaBeta(board chess.Board, depth int, alpha, beta int, color chess.PieceColor, hash uint64) int {
    if ca.stop.Load() {
        return alpha
    }

    originalAlpha := alpha

    // TT lookup now using precomputed hash
    if entry, ok := ca.tt[hash]; ok && entry.depth >= depth {
        switch entry.flag {
        case ttExact:
            return entry.score
        case ttLowerBound:
            if entry.score > alpha {
                alpha = entry.score
            }
        case ttUpperBound:
            if entry.score < beta {
                beta = entry.score
            }
        }
        if alpha >= beta {
            return entry.score
        }
    }

    if depth == 0 {
        return chess.StaticEval(board, color)
    }

    legalMoves := ca.getLegalMovesForColor(board, color)
    if len(legalMoves) == 0 {
        inCheck, _ := isKingInCheck(board, color)
        if inCheck {
            return minScore + 1
        }
        return 0
    }
    legalMoves = ca.orderMoves(board, legalMoves)

    bestMove := legalMoves[0]
    for _, move := range legalMoves {
        next, err := applyMoveToBoard(board, move.FromPos, move.ToPos)
        if err != nil {
            continue
        }

        nextHash := ZobristUpdate(hash, board, move.FromPos, move.ToPos)
        score := -ca.alphaBeta(next, depth-1, -beta, -alpha, chess.OponentColor(color), nextHash)

        if score >= beta {
            ca.tt[hash] = ttEntry{hash: hash, depth: depth, score: beta, flag: ttLowerBound, bestMove: move}
            return beta
        }
        if score > alpha {
            alpha = score
            bestMove = move
        }
    }

    flag := ttExact
    if alpha <= originalAlpha {
        flag = ttUpperBound
    }
    ca.tt[hash] = ttEntry{hash: hash, depth: depth, score: alpha, flag: flag, bestMove: bestMove}

    return alpha
}

func (ca *ChessArbitrator) orderMoves(board chess.Board, moves []chess.Move) []chess.Move {
	type scored struct {
		move  chess.Move
		score int
	}

	sm := make([]scored, 0, len(moves))
	for _, move := range moves {
		score := 0
		target := board.GetPiece(move.ToPos)
		if target != nil {
			attacker := board.GetPiece(move.FromPos)
			score = chess.PieceBaseValue(target.Kind)*10 - chess.PieceBaseValue(attacker.Kind)
		}
		sm = append(sm, scored{move: move, score: score})
	}

	for i := 1; i < len(sm); i++ {
		for j := i; j > 0 && sm[j].score > sm[j-1].score; j-- {
			sm[j], sm[j-1] = sm[j-1], sm[j]
		}
	}

	result := make([]chess.Move, len(sm))
	for i, s := range sm {
		result[i] = s.move
	}

	return result
}

func (ca *ChessArbitrator) getLegalMovesForColor(board chess.Board, color chess.PieceColor) []chess.Move {
	ca.useBoard(board)
	switch color {
	case chess.PieceBlack:
		return ca.legalMovesForPieces(board, ca.blackPiecesPositions)
	case chess.PieceWhite:
		return ca.legalMovesForPieces(board, ca.whitePiecesPositions)
	}

	return make([]chess.Move, 0)
}

func (ca *ChessArbitrator) legalMovesForPieces(board chess.Board, piecePosition []chess.PiecePosition) []chess.Move {
	rows, cols := board.Dimensions()

	var legalMoves []chess.Move
	for _, piecePos := range piecePosition {
		for row := range rows {
			for col := range cols {
				toPos := chess.Position{Row: row, Col: col}
				canMove, _ := ca.CanMovePiece(board, piecePos.Piece, piecePos.Position, toPos)
				if canMove {
					legalMoves = append(legalMoves, chess.Move{FromPos: piecePos.Position, ToPos: toPos})
				}
			}
		}
	}

	return legalMoves
}

func (ca *ChessArbitrator) useBoard(board chess.Board) {
	// clear arrays
	ca.blackPiecesPositions = ca.blackPiecesPositions[:0]
	ca.whitePiecesPositions = ca.whitePiecesPositions[:0]

	rows, cols := board.Dimensions()
	for row := range rows {
		for col := range cols {
			position := chess.Position{Row: row, Col: col}
			piece := board.GetPiece(position)
			if piece == nil {
				continue
			}

			piecePosition := chess.PiecePosition{
				Piece:    *piece,
				Position: position,
			}

			if piece.Color == chess.PieceBlack {
				ca.blackPiecesPositions = append(ca.blackPiecesPositions, piecePosition)
			} else {
				ca.whitePiecesPositions = append(ca.whitePiecesPositions, piecePosition)
			}
		}
	}
}

func canMoveOligarch(board chess.Board, piece chess.Piece, fromPos, toPos chess.Position) (bool, string) {
	hasFriendlyNeighbor := chess.HasAdjacentFriendlyPiece(board, piece.Color, fromPos)
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
	if attackedBySlider(board, squarePosition, color) {
		return true
	}

	knightMoves := []chess.Position{
		{Row: 2, Col: 1}, {Row: 2, Col: -1}, {Row: -2, Col: 1}, {Row: -2, Col: -1},
		{Row: 1, Col: 2}, {Row: 1, Col: -2}, {Row: -1, Col: 2}, {Row: -1, Col: -2},
	}
	if attackedByStepper(board, squarePosition, color, knightMoves, chess.KnightPiece) {
		return true
	}

	kangarooMoves := []chess.Position{
		{Row: 2, Col: 0}, {Row: -2, Col: 0}, {Row: 0, Col: 2}, {Row: 0, Col: -2},
		{Row: 3, Col: 0}, {Row: -3, Col: 0}, {Row: 0, Col: 3}, {Row: 0, Col: -3},
	}
	if attackedByStepper(board, squarePosition, color, kangarooMoves, chess.KangarooPiece) {
		return true
	}

	kingMoves := []chess.Position{
		{Row: 1, Col: 0}, {Row: -1, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: -1},
		{Row: 1, Col: 1}, {Row: 1, Col: -1}, {Row: -1, Col: 1}, {Row: -1, Col: -1},
	}
	if attackedByStepper(board, squarePosition, color, kingMoves, chess.KingPiece) {
		return true
	}

	pawnDir := 1
	if color == chess.PieceWhite {
		pawnDir = -1
	}

	pawnAttacks := []chess.Position{
		{Row: -pawnDir, Col: 1},
		{Row: -pawnDir, Col: -1},
	}
	return attackedByStepper(board, squarePosition, color, pawnAttacks, chess.PawnPiece)
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

func attackedByStepper(board chess.Board, piecePos chess.Position, attackerColor chess.PieceColor, offsets []chess.Position, kind chess.PieceKind) bool {
	for _, offset := range offsets {
		checkPos := chess.Position{Row: piecePos.Row + offset.Row, Col: piecePos.Col + offset.Col}

		if !chess.IsInsideBoard(board, piecePos, checkPos) {
			continue
		}

		piece := board.GetPiece(checkPos)

		if piece != nil && piece.Color == attackerColor && piece.Kind == kind {
			return true
		}
	}

	return false
}

func attackedBySlider(board chess.Board, piecePos chess.Position, attackerColor chess.PieceColor) bool {
	directions := []struct {
		dr, dc   int
		diagonal bool
	}{
		{1, 0, false}, {-1, 0, false}, {0, 1, false}, {0, -1, false}, // orthogonal
		{1, 1, true}, {1, -1, true}, {-1, 1, true}, {-1, -1, true}, // diagonal
	}

	for _, dir := range directions {

		for i := 1; ; i++ {

			currPos := chess.Position{
				Row: piecePos.Row + dir.dr*i,
				Col: piecePos.Col + dir.dc*i,
			}

			if !chess.IsInsideBoard(board, piecePos, currPos) {
				break
			}

			piece := board.GetPiece(currPos)

			if piece == nil {
				continue
			}

			// friendly piece blocks ray
			if piece.Color != attackerColor {
				break
			}

			if isPieceASlider(board, *piece, piecePos, currPos, dir.diagonal) {
				return true
			}

			break
		}
	}

	return false
}

func isPieceASlider(
	board chess.Board,
	piece chess.Piece,
	targetPiecePos chess.Position,
	position chess.Position,
	isDiagonal bool,
) bool {
	switch piece.Kind {
	case chess.QueenPiece:
		return true
	case chess.RookPiece:
		return !isDiagonal
	case chess.BishopPiece:
		return isDiagonal
	case chess.OligarchPiece:
		xDiff, yDiff := distance(targetPiecePos, position)
		// adjacent oligarch behaves like king
		if xDiff <= 1 && yDiff <= 1 {
			return true
		}
		// otherwise it becomes a slider only if activated
		return chess.HasAdjacentFriendlyPiece(board, piece.Color, position)
	default:
		return false
	}
}

func distance(pos1, pos2 chess.Position) (int, int) {
	rowDiff := intAbs(pos1.Row - pos2.Row)
	colDiff := intAbs(pos1.Col - pos2.Col)
	return rowDiff, colDiff
}

func boardHash(board chess.Board, color chess.PieceColor) uint64 {
	const (
		offset uint64 = 14695981039346656037
		prime  uint64 = 1099511628211
	)

	h := offset
	rows, cols := board.Dimensions()

	for row := range rows {
		for col := range cols {
			h ^= uint64(row*100 + col)
			h *= prime

			piece := board.GetPiece(chess.Position{Row: row, Col: col})
			if piece != nil {
				h ^= uint64(piece.Kind)*31 + uint64(piece.Color)*17
				h *= prime
			}
		}
	}

	// include whose turn it is: same position, different turn = different hash
	h ^= uint64(color)
	h *= prime

	return h
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
