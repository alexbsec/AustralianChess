package ws

import (
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	wsTypes "github.com/alexbsec/AustralianChess/backend/internal/ws/ws_types"
)

// Re-export ws_types so callers don't need to import the sub-package.
type (
	Conn       = wsTypes.Conn
	Client     = wsTypes.Client
	RoomClient = wsTypes.RoomClient
	Role       = wsTypes.Role
)

const (
	PlayerOne = wsTypes.PlayerOne
	PlayerTwo = wsTypes.PlayerTwo
	Spectator = wsTypes.Spectator
)

var (
	NewClient     = wsTypes.NewClient
	NewRoomClient = wsTypes.NewRoomClient
)

type IHub interface {
	StartJanitor()
	SetOnRoomEmptyCallback(callback func(roomId string))
	AddClient(roomId, playerId string, conn wsTypes.Conn) (*wsTypes.Client, error)
	AddBot(roomId, botId string, color chess.PieceColor)
	RemoveClient(roomId string, conn wsTypes.Conn)
	Broadcast(roomId string, msg any, resetWarning bool) error
	SetRoomForTest(roomId string, roomClient *wsTypes.RoomClient)
	GetRoomForTest(roomId string) *wsTypes.RoomClient
}
