package ws

import "github.com/gorilla/websocket"

type IHub interface {
	AddClient(roomId, playerId string, conn *websocket.Conn) (*Client, error)
	RemoveClient(roomId string, conn *websocket.Conn)
	Broadcast(roomId string, msg any) error
}
