package ws_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
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


func TestHub_AddClient_FirstClient_CreatesRoom(t *testing.T) {
	hub := ws.NewHub()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := wsMocks.NewMockConn(ctrl)

	client, err := hub.AddClient("room1", "p1", mockConn)

	require.NoError(t, err)
	require.NotNil(t, client)

	room := hub.GetRoomForTest("room1")
	require.NotNil(t, room)
	require.NotNil(t, room.PlayerOne)

	require.Equal(t, "p1", room.PlayerOne.PlayerId)
}


func TestHub_AddClient_SecondClient_AssignsPlayerTwo(t *testing.T) {
	hub := ws.NewHub()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock1 := wsMocks.NewMockConn(ctrl)
	mock2 := wsMocks.NewMockConn(ctrl)

	// Force deterministic behavior (optional but ideal)
	// You may want to mock rand in real refactor

	_, _ = hub.AddClient("room1", "p1", mock1)
	client2, err := hub.AddClient("room1", "p2", mock2)

	require.NoError(t, err)
	require.Equal(t, ws.PlayerTwo, client2.Role)

	room := hub.GetRoomForTest("room1")
	require.NotNil(t, room.PlayerTwo)
	require.NotNil(t, room.PlayerOne.Color)

	// Ensure opposite colors
	if *room.PlayerOne.Color == chess.PieceWhite {
		require.Equal(t, chess.PieceBlack, *room.PlayerTwo.Color)
	} else {
		require.Equal(t, chess.PieceWhite, *room.PlayerTwo.Color)
	}
}


func TestHub_AddClient_ThirdClient_IsSpectator(t *testing.T) {
	hub := ws.NewHub()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock1 := wsMocks.NewMockConn(ctrl)
	mock2 := wsMocks.NewMockConn(ctrl)
	mock3 := wsMocks.NewMockConn(ctrl)

	_, _ = hub.AddClient("room1", "p1", mock1)
	_, _ = hub.AddClient("room1", "p2", mock2)

	client3, err := hub.AddClient("room1", "p3", mock3)

	require.NoError(t, err)
	require.Equal(t, ws.Spectator, client3.Role)

	room := hub.GetRoomForTest("room1")
	require.Contains(t, room.Spectators, mock3)
}


func TestHub_AddClient_Reconnect_PlayerOne(t *testing.T) {
	hub := ws.NewHub()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOld := wsMocks.NewMockConn(ctrl)
	mockNew := wsMocks.NewMockConn(ctrl)

	// Create initial client
	_, _ = hub.AddClient("room1", "p1", mockOld)

	room := hub.GetRoomForTest("room1")
	oldDone := room.PlayerOne.Done

	// Trigger reconnect
	client, err := hub.AddClient("room1", "p1", mockNew)

	require.NoError(t, err)
	require.Equal(t, mockNew, client.Conn)

	// Ensure Done channel was replaced
	require.NotEqual(t, oldDone, client.Done)
}


func TestHub_AddClient_Reconnect_PlayerTwo(t *testing.T) {
	hub := ws.NewHub()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock1 := wsMocks.NewMockConn(ctrl)
	mock2 := wsMocks.NewMockConn(ctrl)
	mockNew := wsMocks.NewMockConn(ctrl)

	_, _ = hub.AddClient("room1", "p1", mock1)
	_, _ = hub.AddClient("room1", "p2", mock2)

	room := hub.GetRoomForTest("room1")
	oldDone := room.PlayerTwo.Done

	client, err := hub.AddClient("room1", "p2", mockNew)

	require.NoError(t, err)
	require.Equal(t, mockNew, client.Conn)
	require.NotEqual(t, oldDone, client.Done)
}
