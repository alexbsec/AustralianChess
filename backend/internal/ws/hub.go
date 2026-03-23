package ws

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
)

type Hub struct {
	mtx         sync.Mutex
	roomClients map[string]*RoomClient
	onRoomEmpty func(roomId string)
}

func NewHub() *Hub {
	return &Hub{
		roomClients: make(map[string]*RoomClient),
	}
}

func (h *Hub) StartJanitor() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			h.cleanupRooms()
		}
	}()
}

func (h *Hub) SetOnRoomEmptyCallback(callback func(roomId string)) {
	h.onRoomEmpty = callback
}

func (h *Hub) AddClient(roomId, playerId string, conn Conn) (*Client, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if _, ok := h.roomClients[roomId]; !ok {
		rc := NewRoomClient()
		rc.LastActive = time.Now()
		h.roomClients[roomId] = rc
	}

	room := h.roomClients[roomId]

	if room.PlayerOne != nil && room.PlayerOne.PlayerId == playerId {
		select {
		case <-room.PlayerOne.Done:
		default:
			close(room.PlayerOne.Done)
		}
		room.PlayerOne.Done = make(chan struct{})
		room.PlayerOne.Conn = conn
		log.Printf("client reconnected as %s", room.PlayerOne.Role)
		return room.PlayerOne, nil
	}

	if room.PlayerTwo != nil && room.PlayerTwo.PlayerId == playerId {
		select {
		case <-room.PlayerTwo.Done:
		default:
			close(room.PlayerTwo.Done)
		}
		room.PlayerTwo.Done = make(chan struct{})
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

func (h *Hub) RemoveClient(roomId string, conn Conn) {
	h.mtx.Lock()

	room, ok := h.roomClients[roomId]
	if !ok {
		h.mtx.Unlock()
		return
	}

	if room.PlayerOne != nil && room.PlayerOne.Conn == conn {
		room.PlayerOne.Conn = nil
		log.Printf("player one disconnected (slot kept for reconnection)")
	}

	if room.PlayerTwo != nil && room.PlayerTwo.Conn == conn {
		room.PlayerTwo.Conn = nil
		log.Printf("player two disconnected (slot kept for reconnection)")
	}

	delete(room.Spectators, conn)

	room.LastActive = time.Now()
	h.mtx.Unlock()
	log.Printf("client disconnected")
}

func (h *Hub) Broadcast(roomId string, message any, resetWarning bool) error {
    h.mtx.Lock()
    room, ok := h.roomClients[roomId]
    if !ok {
        h.mtx.Unlock()
        return nil
    }
    room.LastActive = time.Now()

	if resetWarning {
		room.WarningSent = false
	}

    clients := make([]*Client, 0, 2+len(room.Spectators))
    if room.PlayerOne != nil && room.PlayerOne.Conn != nil {
        clients = append(clients, room.PlayerOne)
    }
    if room.PlayerTwo != nil && room.PlayerTwo.Conn != nil {
        clients = append(clients, room.PlayerTwo)
    }
    for _, c := range room.Spectators {
        clients = append(clients, c)
    }
    h.mtx.Unlock()

    for _, client := range clients {
        if err := client.WriteJSON(message); err != nil {
            log.Printf("broadcast failed: %v", err)
        }
    }
    return nil
}

func (h *Hub) SetRoomForTest(roomId string, roomClient *RoomClient) {
	h.addRoomForTest(roomId, roomClient)
}

func (h *Hub) GetRoomForTest(roomId string) *RoomClient {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	room, ok := h.roomClients[roomId]
	if !ok {
		return nil
	}

	return room
}

func (h *Hub) cleanupRooms() {
	h.mtx.Lock()

	var deleted []string
	var warned []string

	for id, room := range h.roomClients {
		inactiveSince := time.Since(room.LastActive)
		maxTimeMin := 7 * time.Minute
		warnAt := 6 * time.Minute + 30 * time.Second

		if inactiveSince >= maxTimeMin {
			delete(h.roomClients, id)
			deleted = append(deleted, id)
		} else if inactiveSince >= warnAt && !room.WarningSent {
			room.WarningSent = true
			warned = append(warned, id)
		}
	}
	h.mtx.Unlock()

	for _, id := range warned {
		h.Broadcast(id, map[string]any{
			"type": "inactive_warning",
			"seconds": 30,
		}, false)
	}

	if h.onRoomEmpty == nil {
		return
	}

	for _, id := range deleted {
		h.onRoomEmpty(id)
		log.Printf("room %s cleaned up by janitor", id)
	}
}

func (h *Hub) addRoomForTest(roomId string, room *RoomClient) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	h.roomClients[roomId] = room
}
