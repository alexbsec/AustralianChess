package ws

import (
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
}

type RoomClient struct {
	PlayerOne  *Client
	PlayerTwo  *Client
	Spectators map[*websocket.Conn]*Client
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
