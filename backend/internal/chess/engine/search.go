package engine

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type ttEntry struct {
	hash     uint64
	depth    int
	score    int
	flag     int
	bestMove chess.Move
}
