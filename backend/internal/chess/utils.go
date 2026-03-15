package chess

func PieceColorString(color PieceColor) string {
	switch color {
	case PieceWhite:
		return "white"
	case PieceBlack:
		return "black"
	default:
		return "unknown"
	}
}
