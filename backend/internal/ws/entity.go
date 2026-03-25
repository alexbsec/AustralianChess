package ws

import (
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
)

type PlayBotDTO struct {
	RoomId          string           `json:"roomId"`
	Difficulty      bot.Difficulty   `json:"difficulty"`
	PlayerPlayingAs chess.PieceColor `json:"player_playing_as"`
}
