package ws_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	wsMocks "github.com/alexbsec/AustralianChess/backend/internal/ws/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHub_Broadcast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := wsMocks.NewMockConn(ctrl)

	hub := ws.NewHub()

	client := ws.NewClient("room1", nil).WithPlayerId("p1").WithRole(ws.PlayerOne)
	client.Conn = mockConn

	room := ws.NewRoomClient()
	room.PlayerOne = client

	hub.SetRoomForTest("room1", room) // use helper

	mockConn.EXPECT().
		WriteJSON(gomock.Any()).
		Return(nil)

	err := hub.Broadcast("room1", "hello", false)
	require.NoError(t, err)
}

func TestHub_Broadcast_MultipleClients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock1 := wsMocks.NewMockConn(ctrl)
	mock2 := wsMocks.NewMockConn(ctrl)

	hub := ws.NewHub()

	c1 := ws.NewClient("room1", nil).WithRole(ws.PlayerOne)
	c1.Conn = mock1

	c2 := ws.NewClient("room1", nil).WithRole(ws.PlayerTwo)
	c2.Conn = mock2

	room := ws.NewRoomClient()
	room.PlayerOne = c1
	room.PlayerTwo = c2

	hub.SetRoomForTest("room1", room)

	mock1.EXPECT().WriteJSON(gomock.Any()).Return(nil)
	mock2.EXPECT().WriteJSON(gomock.Any()).Return(nil)

	err := hub.Broadcast("room1", "msg", false)
	require.NoError(t, err)
}

func TestHub_Broadcast_UnknownRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := ws.NewHub()

	mock := wsMocks.NewMockConn(ctrl)
	mock.EXPECT().WriteJSON(gomock.Any()).Times(0)

	err := hub.Broadcast("unknown", "msg", false)
	require.NoError(t, err)
}


func TestHub_Broadcast_WriteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := wsMocks.NewMockConn(ctrl)

	hub := ws.NewHub()

	client := ws.NewClient("room1", nil).WithRole(ws.PlayerOne)
	client.Conn = mock

	room := ws.NewRoomClient()
	room.PlayerOne = client

	hub.SetRoomForTest("room1", room)

	mock.EXPECT().
		WriteJSON(gomock.Any()).
		Return(assert.AnError)

	err := hub.Broadcast("room1", "msg", false)
	require.NoError(t, err) // should not fail
}


func TestHub_RemoveClient_PlayerOne(t *testing.T) {
	hub := ws.NewHub()

	mock := wsMocks.NewMockConn(gomock.NewController(t))

	client := ws.NewClient("room1", mock).WithRole(ws.PlayerOne)

	room := ws.NewRoomClient()
	room.PlayerOne = client

	hub.SetRoomForTest("room1", room)

	hub.RemoveClient("room1", mock)

	require.Nil(t, room.PlayerOne.Conn)
}
