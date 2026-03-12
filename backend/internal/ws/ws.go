package ws

import "github.com/gorilla/websocket"

type IHub interface {
	AddClient(roomId string, conn *websocket.Conn)
	RemoveClient(roomId string, conn *websocket.Conn)
	Broadcast(roomId string, msg any) error
}
