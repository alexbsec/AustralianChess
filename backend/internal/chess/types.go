package chess

type Move struct {
	FromPos Position
	ToPos   Position
}

type PiecePosition struct {
	Piece    Piece
	Position Position
}

type PieceMove struct {
	Piece *Piece
	Move  Move
}
