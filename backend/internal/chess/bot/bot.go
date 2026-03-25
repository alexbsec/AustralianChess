package bot

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type Bot interface {
	Start()
	Stop()
	NotifyState(state chess.GameState)
}
