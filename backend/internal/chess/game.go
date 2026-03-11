package chess

type GameResult int

const (
	ResultCheckmate GameResult = iota
	ResultStalemate
)

type GameState struct {
	Board     Board       `json:"board"`
	Turn      PieceColor  `json:"turn"`
	Result    *GameResult `json:"result"`
	EndReason *string     `json:"endReason"`
}

func NewGame(board Board) *GameState {
	return &GameState{
		Board:     board,
		Turn:      PieceWhite,
		Result:    nil,
		EndReason: nil,
	}
}
