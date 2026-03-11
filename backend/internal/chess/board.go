package chess

type SquareColor int

const (
	SquareColorLight SquareColor = iota
	SquareColorDark
)

const rows = 12
const cols = 12

type Square struct {
	Color SquareColor `json:"color"`
	Piece *Piece      `json:"piece"`
}

type Board struct {
	Data [][]Square `json:"data"`
}

func NewInitialBoard() Board {
	board := Board{
		Data: make([][]Square, rows),
	}

	for r := range rows {
		board.Data[r] = make([]Square, cols)

		for c := range cols {
			color := SquareColorLight
			if (r+c)%2 == 1 {
				color = SquareColorDark
			}

			board.Data[r][c] = Square{
				Color: color,
				Piece: nil,
			}
		}
	}

	addPiecesToBoard(&board)
	return board
}

func addPiecesToBoard(board *Board) {
	backRank := [cols]PieceKind{
		RookPiece,
		OligarchPiece,
		KnightPiece,
		KangarooPiece,
		BishopPiece,
		QueenPiece,
		KingPiece,
		BishopPiece,
		KangarooPiece,
		KnightPiece,
		OligarchPiece,
		RookPiece,
	}

	for col, pieceName := range backRank {
		board.Data[0][col].Piece = NewPiece(pieceName, PieceBlack)
		board.Data[11][col].Piece = NewPiece(pieceName, PieceWhite)
	}

	// Pawns
	for col := range cols {
		board.Data[1][col].Piece = NewPiece(PawnPiece, PieceBlack)
		board.Data[10][col].Piece = NewPiece(PawnPiece, PieceWhite)
	}
}
