package chess

import "fmt"

type SquareColor int

const (
	SquareColorLight SquareColor = iota
	SquareColorDark
)

const rows = 12
const cols = 12

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

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

func (b Board) GetPiece(position Position) *Piece {
	square := b.Data[position.Row][position.Col]
	return square.Piece
}

func (b *Board) MovePiece(fromPos Position, toPos Position) error {
	piece := b.GetPiece(fromPos)
	if piece == nil {
		return fmt.Errorf("no piece to move at position: %v", fromPos)
	}

	b.Data[fromPos.Row][fromPos.Col].Piece = nil
	b.Data[toPos.Row][toPos.Col].Piece = piece
	return nil
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
