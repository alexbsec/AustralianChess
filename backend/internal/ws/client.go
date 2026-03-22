package ws

import (
	"sync"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/gorilla/websocket"
)

type Role string

const (
	PlayerOne Role = "Player1"
	PlayerTwo Role = "Player2"
	Spectator Role = "Spectator"
)

type Client struct {
	Conn     *websocket.Conn
	RoomId   string
	PlayerId string
	Color    *chess.PieceColor
	Role     Role
	WriteMtx sync.Mutex
	Done     chan struct{}
}

type RoomClient struct {
	PlayerOne  *Client
	PlayerTwo  *Client
	Spectators map[*websocket.Conn]*Client
	LastActive time.Time
}

func NewRoomClient() *RoomClient {
	return &RoomClient{
		Spectators: make(map[*websocket.Conn]*Client),
	}
}

func NewClient(roomId string, conn *websocket.Conn) *Client {
	return &Client{
		Conn:   conn,
		RoomId: roomId,
		Done: make(chan struct{}),
	}
}

func (c *Client) WithPlayerId(id string) *Client {
	c.PlayerId = id
	return c
}

func (c *Client) WithColor(color chess.PieceColor) *Client {
	c.Color = &color
	return c
}

func (c *Client) WithRole(role Role) *Client {
	c.Role = role
	return c
}

func (c *Client) WriteJSON(v any) error {
	c.WriteMtx.Lock()
	defer c.WriteMtx.Unlock()
	return c.Conn.WriteJSON(v)
}

func (c *Client) WritePing() error {
	c.WriteMtx.Lock()
	defer c.WriteMtx.Unlock()
	return c.Conn.WriteMessage(websocket.PingMessage, nil)
}
