package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mtx   sync.RWMutex
	rooms map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) AddClient(roomId string, conn *websocket.Conn) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if _, ok := h.rooms[roomId]; !ok {
		h.rooms[roomId] = make(map[*websocket.Conn]struct{})
	}

	h.rooms[roomId][conn] = struct{}{}
	log.Printf("client connected")
}

func (h *Hub) RemoveClient(roomId string, conn *websocket.Conn) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	clients, ok := h.rooms[roomId]
	if !ok {
		return
	}

	delete(clients, conn)
	if len(clients) == 0 {
		delete(h.rooms, roomId)
	}

	log.Printf("client disconnected")
}

func (h *Hub) Broadcast(roomId string, message any) error {
	return nil
}
