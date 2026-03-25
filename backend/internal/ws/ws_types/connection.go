package ws_types

import (
	"time"
)

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
