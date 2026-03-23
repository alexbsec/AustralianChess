package ws

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	contractsMocks "github.com/alexbsec/AustralianChess/backend/internal/contracts/mocks"
	parserMocks "github.com/alexbsec/AustralianChess/backend/internal/ws/parser/mocks"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	roomsMocks "github.com/alexbsec/AustralianChess/backend/rooms/mocks"
	"github.com/alexbsec/AustralianChess/backend/users"
	usersMocks "github.com/alexbsec/AustralianChess/backend/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// testHub: wraps a real Hub with hooks needed by loop tests.
//
// ws/mocks cannot be imported from package ws (it would be a circular
// dependency), so we use a lightweight wrapper instead.
// ---------------------------------------------------------------------------

type testHub struct {
	*Hub
	mu              sync.Mutex
	broadcastCalls  int
	failBroadcastAt int // 0 = never; N = fail on the Nth call
	onRemoveClient  func()
}

func newTestHub(failAt int, onRemove func()) *testHub {
	return &testHub{
		Hub:             NewHub(),
		failBroadcastAt: failAt,
		onRemoveClient:  onRemove,
	}
}

// Broadcast overrides the real Hub so we can inject a failure on demand.
func (h *testHub) Broadcast(roomId string, msg any, resetWarning bool) error {
	h.mu.Lock()
	h.broadcastCalls++
	n := h.broadcastCalls
	h.mu.Unlock()
	if h.failBroadcastAt > 0 && n == h.failBroadcastAt {
		return errors.New("broadcast failed")
	}
	return h.Hub.Broadcast(roomId, msg, resetWarning)
}

// RemoveClient calls the real implementation then fires the optional hook so
// tests can synchronise on handler-goroutine completion.
func (h *testHub) RemoveClient(roomId string, conn Conn) {
	h.Hub.RemoveClient(roomId, conn)
	if h.onRemoveClient != nil {
		h.onRemoveClient()
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newPrivateRouter(t *testing.T, h *Handler) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ws/:id", func(c *gin.Context) {
		c.Set("sessionId", int64(1))
		c.Set("userId", int64(42))
		h.HandleRoom(c)
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func privateWsURL(srv *httptest.Server, roomId string) string {
	return "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/" + roomId
}

// dialAndDrainN dials the given URL and discards the first n messages, leaving
// the connection in the "inside the loop" state.
func dialAndDrainN(t *testing.T, url string, n int) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	for i := 0; i < n; i++ {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _, readErr := conn.ReadMessage()
		require.NoError(t, readErr)
	}
	return conn
}

// waitHandlerDone blocks until wg reaches zero or the 2-second deadline fires.
func waitHandlerDone(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for handler goroutine to finish")
	}
}

func newPrivateGameState() chess.GameState {
	return *chess.NewGame(chess.NewInitialBoard())
}

// loopTestDeps groups the objects needed by every player-loop test.
type loopTestDeps struct {
	hub        *testHub
	roomSvc    *roomsMocks.MockIService
	mockParser *parserMocks.MockIParser
	handler    *Handler
}

// setupPlayerLoopTest wires up all mocks for a PlayerOne connecting to room-1.
// failBroadcastAt controls on which Broadcast call testHub should return an
// error (0 = never).  The returned WaitGroup is Done when RemoveClient fires.
func setupPlayerLoopTest(
	t *testing.T,
	ctrl *gomock.Controller,
	failBroadcastAt int,
) (*loopTestDeps, *sync.WaitGroup) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(1)

	hub := newTestHub(failBroadcastAt, func() { wg.Done() })
	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	mockParser := parserMocks.NewMockIParser(ctrl)

	gs := newPrivateGameState()
	joinResult := contracts.PlayerJoinedResult{
		RoomId:   "room-1",
		PlayerId: "p1",
		Success:  true,
	}

	userSvc.EXPECT().
		User(gomock.Any(), int64(42)).
		Return(&users.User{PlayerId: "p1"}, nil)
	roomSvc.EXPECT().
		FetchRoom(gomock.Any(), "room-1").
		Return(&rooms.Room{Id: "room-1", GameState: &gs}, nil)
	// Real Hub randomises color, so match any value.
	roomSvc.EXPECT().
		UpdatePlayerJoined(gomock.Any(), "room-1", "p1", gomock.Any()).
		Return(joinResult, nil)

	handler := NewHandler(roomSvc, userSvc, hub)
	handler.wsParser = mockParser // only possible from within package ws

	return &loopTestDeps{
		hub:        hub,
		roomSvc:    roomSvc,
		mockParser: mockParser,
		handler:    handler,
	}, &wg
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestHandleRoom_Loop_SpectatorSendsMessage_GetsError verifies that the loop
// immediately sends an error and continues when a spectator sends a command.
func TestHandleRoom_Loop_SpectatorSendsMessage_GetsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var wg sync.WaitGroup
	wg.Add(1)

	hub := newTestHub(0, func() { wg.Done() })

	// Pre-populate the room so that "p1" is assigned the Spectator role.
	rc := NewRoomClient()
	rc.PlayerOne = &Client{PlayerId: "p2", Role: PlayerOne, Done: make(chan struct{})}
	rc.PlayerTwo = &Client{PlayerId: "p3", Role: PlayerTwo, Done: make(chan struct{})}
	hub.SetRoomForTest("room-1", rc)

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)

	gs := newPrivateGameState()
	userSvc.EXPECT().User(gomock.Any(), int64(42)).Return(&users.User{PlayerId: "p1"}, nil)
	roomSvc.EXPECT().FetchRoom(gomock.Any(), "room-1").Return(&rooms.Room{Id: "room-1", GameState: &gs}, nil)

	handler := NewHandler(roomSvc, userSvc, hub)
	srv := newPrivateRouter(t, handler)

	// Spectator receives only the initial game_state (no join broadcast).
	conn := dialAndDrainN(t, privateWsURL(srv, "room-1"), 1)

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"move"}`)))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(msg, &payload))
	assert.Equal(t, "error", payload["type"])
	assert.Equal(t, "spectators cannot send commands", payload["message"])

	_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
	waitHandlerDone(t, &wg)
}

// TestHandleRoom_Loop_ParseFails_SendsError verifies that a parse failure
// causes the handler to write an error back to the client and continue looping.
func TestHandleRoom_Loop_ParseFails_SendsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deps, wg := setupPlayerLoopTest(t, ctrl, 0)

	deps.mockParser.EXPECT().
		ParseMessage("room-1", gomock.Any()).
		Return(nil, errors.New("bad message"))

	srv := newPrivateRouter(t, deps.handler)
	// Player receives game_state (msg 1) and join-result broadcast (msg 2).
	conn := dialAndDrainN(t, privateWsURL(srv, "room-1"), 2)

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"move"}`)))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(msg, &payload))
	assert.Equal(t, "error", payload["type"])
	assert.Equal(t, "bad message", payload["message"])

	_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
	waitHandlerDone(t, wg)
}

// TestHandleRoom_Loop_ExecuteCommandFails_SendsError verifies that an
// ExecuteCommand failure causes the handler to send an error to the client.
func TestHandleRoom_Loop_ExecuteCommandFails_SendsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deps, wg := setupPlayerLoopTest(t, ctrl, 0)

	mockCmd := contractsMocks.NewMockCommand(ctrl)
	deps.mockParser.EXPECT().
		ParseMessage("room-1", gomock.Any()).
		Return(mockCmd, nil)
	mockCmd.EXPECT().SetPlayerId("p1").Return(mockCmd)
	mockCmd.EXPECT().SetColor(gomock.Any()).Return(mockCmd)
	deps.roomSvc.EXPECT().
		ExecuteCommand(gomock.Any(), mockCmd).
		Return(nil, errors.New("illegal move"))

	srv := newPrivateRouter(t, deps.handler)
	conn := dialAndDrainN(t, privateWsURL(srv, "room-1"), 2)

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"move"}`)))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(msg, &payload))
	assert.Equal(t, "error", payload["type"])
	assert.Equal(t, "illegal move", payload["message"])

	_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
	waitHandlerDone(t, wg)
}

// TestHandleRoom_Loop_BroadcastFails_ConnectionClosed verifies that a
// Broadcast failure after a successful command causes the handler to close the
// connection.
//
// failBroadcastAt=2: the join broadcast (call 1) succeeds, the command
// broadcast (call 2) returns an error so the loop exits.
func TestHandleRoom_Loop_BroadcastFails_ConnectionClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deps, wg := setupPlayerLoopTest(t, ctrl, 2)

	mockCmd := contractsMocks.NewMockCommand(ctrl)
	mockResult := contractsMocks.NewMockResult(ctrl)

	deps.mockParser.EXPECT().
		ParseMessage("room-1", gomock.Any()).
		Return(mockCmd, nil)
	mockCmd.EXPECT().SetPlayerId("p1").Return(mockCmd)
	mockCmd.EXPECT().SetColor(gomock.Any()).Return(mockCmd)
	deps.roomSvc.EXPECT().
		ExecuteCommand(gomock.Any(), mockCmd).
		Return(mockResult, nil)

	srv := newPrivateRouter(t, deps.handler)
	conn := dialAndDrainN(t, privateWsURL(srv, "room-1"), 2)

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"move"}`)))

	// The server closes the connection after the failed broadcast.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, readErr := conn.ReadMessage()
	assert.Error(t, readErr)

	conn.Close()
	waitHandlerDone(t, wg)
}

// TestHandleRoom_Loop_SuccessfulCommand_Broadcasts verifies the happy path:
// a valid command is parsed, executed, and its result is broadcast to the room.
func TestHandleRoom_Loop_SuccessfulCommand_Broadcasts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deps, wg := setupPlayerLoopTest(t, ctrl, 0)

	mockCmd := contractsMocks.NewMockCommand(ctrl)
	moveResult := contracts.MoveResult{RoomId: "room-1", Moved: true}

	deps.mockParser.EXPECT().
		ParseMessage("room-1", gomock.Any()).
		Return(mockCmd, nil)
	mockCmd.EXPECT().SetPlayerId("p1").Return(mockCmd)
	mockCmd.EXPECT().SetColor(gomock.Any()).Return(mockCmd)
	deps.roomSvc.EXPECT().
		ExecuteCommand(gomock.Any(), mockCmd).
		Return(moveResult, nil)

	srv := newPrivateRouter(t, deps.handler)
	conn := dialAndDrainN(t, privateWsURL(srv, "room-1"), 2)

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"move"}`)))

	// The real Hub.Broadcast writes the move result back to the client.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(msg, &payload))
	assert.Equal(t, "room-1", payload["room_id"])
	assert.Equal(t, true, payload["moved"])

	_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
	waitHandlerDone(t, wg)
}
