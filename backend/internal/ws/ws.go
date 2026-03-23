package ws

import (
	"time"
)

type IHub interface {
	StartJanitor()
	SetOnRoomEmptyCallback(callback func(roomId string))
	AddClient(roomId, playerId string, conn Conn) (*Client, error)
	RemoveClient(roomId string, conn Conn)
	Broadcast(roomId string, msg any, resetWarning bool) error
	SetRoomForTest(roomId string, roomClient *RoomClient)
}

type Conn interface {
	// writing
	WriteJSON(v any) error
	WriteMessage(messageType int, data []byte) error
	WriteControl(messageType int, data []byte, deadline time.Time) error

	// reading
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error

	// lifecycle
	Close() error
}
