package rooms

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type CreateRoomResponse struct {
	RoomId string `json:"roomId"`
}

type DisplayRoomResponse struct {
	RoomId    string          `json:"roomId"`
	Status    string          `json:"status"`
	GameState chess.GameState `json:"gameState"`
}
