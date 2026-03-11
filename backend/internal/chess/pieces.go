package chess

type PieceColor int
type PieceKind int

const (
	PieceWhite PieceColor = iota
	PieceBlack
)

const (
	PawnPiece PieceKind = iota
	RookPiece
	OligarchPiece
	KnightPiece
	KangarooPiece
	BishopPiece
	KingPiece
	QueenPiece
)

type Piece struct {
	Kind  PieceKind  `json:"kind"`
	Color PieceColor `json:"color"`
}

func NewPiece(kind PieceKind, color PieceColor) *Piece {
	return &Piece{
		Kind:  kind,
		Color: color,
	}
}
