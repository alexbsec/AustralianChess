package ws

import (
	"log"
	"math/rand"
	"sync"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mtx         sync.RWMutex
	roomClients map[string]*RoomClient
}

func NewHub() *Hub {
	return &Hub{
		roomClients: make(map[string]*RoomClient),
	}
}

func (h *Hub) AddClient(roomId, playerId string, conn *websocket.Conn) (*Client, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if _, ok := h.roomClients[roomId]; !ok {
		h.roomClients[roomId] = NewRoomClient()
	}

	room := h.roomClients[roomId]

	if room.PlayerOne != nil && room.PlayerOne.PlayerId == playerId {
		room.PlayerOne.Conn = conn
		log.Printf("client reconnected as %s", room.PlayerOne.Role)
		return room.PlayerOne, nil
	}

	if room.PlayerTwo != nil && room.PlayerTwo.PlayerId == playerId {
		room.PlayerTwo.Conn = conn
		log.Printf("client reconnected as %s", room.PlayerTwo.Role)
		return room.PlayerTwo, nil
	}

	var client *Client

	if room.PlayerOne == nil {
		var color chess.PieceColor
		if rand.Intn(2) == 0 {
			color = chess.PieceWhite
		} else {
			color = chess.PieceBlack
		}

		client = NewClient(roomId, conn).
			WithPlayerId(playerId).
			WithRole(PlayerOne).
			WithColor(color)

		room.PlayerOne = client
	} else if room.PlayerTwo == nil {
		var color chess.PieceColor
		if *room.PlayerOne.Color == chess.PieceWhite {
			color = chess.PieceBlack
		} else {
			color = chess.PieceWhite
		}

		client = NewClient(roomId, conn).
			WithPlayerId(playerId).
			WithRole(PlayerTwo).
			WithColor(color)

		room.PlayerTwo = client
	} else {
		client = NewClient(roomId, conn).
			WithPlayerId(playerId).
			WithRole(Spectator)

		room.Spectators[conn] = client
	}

	log.Printf("client connected as %s", client.Role)
	return client, nil
}

func (h *Hub) RemoveClient(roomId string, conn *websocket.Conn) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	room, ok := h.roomClients[roomId]
	if !ok {
		return
	}

	if room.PlayerOne != nil && room.PlayerOne.Conn == conn {
		room.PlayerOne.Conn = nil
		log.Printf("player one disconnected")
	}

	if room.PlayerTwo != nil && room.PlayerTwo.Conn == conn {
		room.PlayerTwo.Conn = nil
		log.Printf("player two disconnected")
	}

	delete(room.Spectators, conn)

	isPlayerOneGone := room.PlayerOne == nil || room.PlayerOne.Conn == nil
	isPlayerTwoGone := room.PlayerTwo == nil || room.PlayerTwo.Conn == nil

	if isPlayerOneGone && isPlayerTwoGone && len(room.Spectators) == 0 {
		delete(h.roomClients, roomId)
	}

	log.Printf("client disconnected")
}

func (h *Hub) Broadcast(roomId string, message any) error {
	h.mtx.RLock()

	room, ok := h.roomClients[roomId]
	if !ok {
		h.mtx.RUnlock()
		return nil
	}

	clients := make([]*websocket.Conn, 0, 2+len(room.Spectators))

	if room.PlayerOne != nil && room.PlayerOne.Conn != nil {
		clients = append(clients, room.PlayerOne.Conn)
	}

	if room.PlayerTwo != nil && room.PlayerTwo.Conn != nil {
		clients = append(clients, room.PlayerTwo.Conn)
	}

	for conn := range room.Spectators {
		clients = append(clients, conn)
	}

	h.mtx.RUnlock()

	for _, conn := range clients {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("broadcast failed: %v", err)
		}
	}

	return nil
}
