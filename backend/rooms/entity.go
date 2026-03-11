package rooms

import (
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

type Room struct {
	Id        string
	GameState *chess.GameState
	WhiteId   string
	BlackId   string
	CreatedAt time.Time
}
