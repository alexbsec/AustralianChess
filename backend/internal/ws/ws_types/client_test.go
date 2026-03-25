package ws_types

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func newTestWebSocket(t *testing.T) (*websocket.Conn, func()) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)

		// Keep connection alive
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	}))
	t.Cleanup(server.Close)

	wsURL := "ws" + server.URL[len("http"):]

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	return conn, func() {
		conn.Close()
	}
}

func TestNewClient(t *testing.T) {
	conn, cleanup := newTestWebSocket(t)
	defer cleanup()

	client := NewClient("room-1", conn)

	require.Equal(t, "room-1", client.RoomId)
	require.NotNil(t, client.Done)
	require.NotNil(t, client.Conn)
}

func TestClient_FluentMethods(t *testing.T) {
	conn, cleanup := newTestWebSocket(t)
	defer cleanup()

	client := NewClient("room-1", conn)

	client = client.
		WithPlayerId("player-1").
		WithRole(PlayerOne).
		WithColor(chess.PieceWhite)

	require.Equal(t, "player-1", client.PlayerId)
	require.Equal(t, PlayerOne, client.Role)
	require.NotNil(t, client.Color)
	require.Equal(t, chess.PieceWhite, *client.Color)
}

func TestClient_WriteJSON(t *testing.T) {
	conn, cleanup := newTestWebSocket(t)
	defer cleanup()

	client := &Client{
		Conn:     conn,
		WriteMtx: sync.Mutex{},
	}

	// Start a reader so the connection stays healthy
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	payload := map[string]string{"msg": "hello"}

	err := client.WriteJSON(payload)
	require.NoError(t, err)
}


func TestClient_WritePing(t *testing.T) {
	conn, cleanup := newTestWebSocket(t)
	defer cleanup()

	client := &Client{
		Conn:     conn,
		WriteMtx: sync.Mutex{},
	}

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	err := client.WritePing()
	require.NoError(t, err)
}


func TestNewRoomClient(t *testing.T) {
	rc := NewRoomClient()

	require.NotNil(t, rc)
	require.NotNil(t, rc.Spectators)
	require.Len(t, rc.Spectators, 0)
	require.True(t, rc.LastActive.IsZero())
	require.False(t, rc.WarningSent)
}


func TestRoomClient_Spectators(t *testing.T) {
	conn, cleanup := newTestWebSocket(t)
	defer cleanup()

	client := NewClient("room-1", conn)

	rc := NewRoomClient()
	rc.Spectators[conn] = client

	require.Len(t, rc.Spectators, 1)
	require.Equal(t, client, rc.Spectators[conn])
}


func TestRoomClient_PlayerAssignment(t *testing.T) {
	rc := NewRoomClient()

	conn1, cleanup1 := newTestWebSocket(t)
	defer cleanup1()

	conn2, cleanup2 := newTestWebSocket(t)
	defer cleanup2()

	p1 := NewClient("room", conn1).WithRole(PlayerOne)
	p2 := NewClient("room", conn2).WithRole(PlayerTwo)

	rc.PlayerOne = p1
	rc.PlayerTwo = p2

	require.Equal(t, p1, rc.PlayerOne)
	require.Equal(t, p2, rc.PlayerTwo)
}


func TestRoomClient_State(t *testing.T) {
	rc := NewRoomClient()

	now := time.Now()
	rc.LastActive = now
	rc.WarningSent = true

	require.Equal(t, now, rc.LastActive)
	require.True(t, rc.WarningSent)
}
