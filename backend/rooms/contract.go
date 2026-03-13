package rooms

import "github.com/alexbsec/AustralianChess/backend/internal/chess"

type CreateRoomResponse struct {
	RoomId string `json:"room_id"`
}

type DisplayRoomResponse struct {
	RoomId    string          `json:"room_id"`
	Status    string          `json:"status"`
	GameState chess.GameState `json:"game_state"`
}
