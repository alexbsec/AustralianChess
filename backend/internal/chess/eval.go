package chess

var pieceBaseValueMap = map[PieceKind]int{
    PawnPiece:     100,
    KnightPiece:   320,
    KangarooPiece: 280,
    BishopPiece:   380,
    RookPiece:     550,
    OligarchPiece: 300, 
    QueenPiece:    1000,
    KingPiece:     20000,
}

func EvaluateBoard(board Board) int {
    score := 0
    rows, cols := board.Dimensions()
    for row := range rows {
        for col := range cols {
            pos := Position{Row: row, Col: col}
            piece := board.GetPiece(pos)
            if piece == nil {
                continue
            }
            value := PieceValue(board, *piece, pos, rows, cols)
            if piece.Color == PieceWhite {
                score += value
            } else {
                score -= value
            }
        }
    }
    return score
}

func StaticEval(board Board, color PieceColor) int {
	raw := EvaluateBoard(board)
	if color == PieceWhite {
		return raw
	}

	return -raw
}

func IsInsideBoard(board Board, fromPos, toPos Position) bool {
	rows, cols := board.Dimensions()
	if fromPos.Row < 0 || fromPos.Row >= rows || fromPos.Col < 0 || fromPos.Col >= cols ||
		toPos.Row < 0 || toPos.Row >= rows || toPos.Col < 0 || toPos.Col >= cols {
		return false
	}

	return true
}

func HasAdjacentFriendlyPiece(board Board, color PieceColor, position Position) bool {
	for rowOffset := -1; rowOffset <= 1; rowOffset++ {
		for colOffset := -1; colOffset <= 1; colOffset++ {
			if rowOffset == 0 && colOffset == 0 {
				continue
			}

			nextRow := position.Row + rowOffset
			nextCol := position.Col + colOffset
			nextPos := Position{Row: nextRow, Col: nextCol}

			if !IsInsideBoard(board, position, nextPos) {
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

func PieceBaseValue(pieceKind PieceKind) int {
	return pieceBaseValueMap[pieceKind]
}

func PieceValue(board Board, piece Piece, position Position, rows, cols int) int {
	base := pieceBaseValueMap[piece.Kind]

	switch piece.Kind {
	case PawnPiece:
		return base + pawnBonus(piece, position)
	case KnightPiece:
		return base + mobilityBonus(position, rows, cols)
	case KangarooPiece:
		return base + kangarooBonus(board, piece, position, rows, cols)
	case BishopPiece:
		return base + bishopBonus(board, position)
	case RookPiece:
		return base + rookBonus(board, piece, position, rows)
	case OligarchPiece:
		return base + oligarchBonus(board, piece, position)
	default:
		return base
	}
}

func pawnBonus(piece Piece, pos Position) int {
    if piece.Color == PieceWhite {
        return max(0, 10-pos.Row) * 15
    }
    return max(0, pos.Row-1) * 15
}

func mobilityBonus(pos Position, rows, cols int) int {
	// penalizes edge squares for jumping pieces
	distFromEdge := min4(pos.Row, pos.Col, rows - 1 - pos.Row, cols - 1 - pos.Col)
	if distFromEdge == 0 {
		return -20
	}

	if distFromEdge >= 3 {
		return 15
	}

	return 0
}

// kangarooBonus counts reachable jump squares as a mobility bonus.
func kangarooBonus(board Board, piece Piece, pos Position, rows, cols int) int {
    targets := []Position{
        {Row: pos.Row + 2, Col: pos.Col},
        {Row: pos.Row - 2, Col: pos.Col},
        {Row: pos.Row, Col: pos.Col + 2},
        {Row: pos.Row, Col: pos.Col - 2},
        {Row: pos.Row + 3, Col: pos.Col},
        {Row: pos.Row - 3, Col: pos.Col},
        {Row: pos.Row, Col: pos.Col + 3},
        {Row: pos.Row, Col: pos.Col - 3},
    }
    reachable := 0
    for _, t := range targets {
        if t.Row < 0 || t.Row >= rows || t.Col < 0 || t.Col >= cols {
            continue
        }
        target := board.GetPiece(t)
        if target == nil || target.Color != piece.Color {
            reachable++
        }
    }
    if reachable >= 4 {
        return 20
    }
    return 0
}

// bishopBonus rewards open diagonals.
func bishopBonus(board Board, pos Position) int {
    bonus := 0
    diagonals := [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
    for _, d := range diagonals {
        open := 0
        r, c := pos.Row+d[0], pos.Col+d[1]
        for r >= 0 && r < 12 && c >= 0 && c < 12 {
            if board.GetPiece(Position{Row: r, Col: c}) != nil {
                break
            }
            open++
            r += d[0]
            c += d[1]
        }
        if open >= 6 {
            bonus += 30
        }
    }
    return bonus
}

func rookBonus(board Board, piece Piece, pos Position, rows int) int {
    bonus := 0

    // Open file bonus
    openFile := true
    for row := range rows {
        if row == pos.Row {
            continue
        }
        p := board.GetPiece(Position{Row: row, Col: pos.Col})
        if p != nil && p.Kind == PawnPiece {
            openFile = false
            break
        }
    }
    if openFile {
        bonus += 30
    }

    // 7th/8th rank bonus (opponent's territory)
    if piece.Color == PieceWhite && pos.Row <= 2 {
        bonus += 20
    } else if piece.Color == PieceBlack && pos.Row >= rows-3 {
        bonus += 20
    }

    return bonus
}

func oligarchBonus(board Board, piece Piece, pos Position) int {
    if HasAdjacentFriendlyPiece(board, piece.Color, pos) {
        return 600
    }
    return 0
}

func min4(a, b, c, d int) int {
    return min(a, min(b, min(c, d)))
}
