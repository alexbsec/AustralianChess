package ws_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	ws "github.com/alexbsec/AustralianChess/backend/internal/ws"
	wsMocks "github.com/alexbsec/AustralianChess/backend/internal/ws/mocks"
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

func init() {
	gin.SetMode(gin.TestMode)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newRouter(t *testing.T, handler *ws.Handler) *httptest.Server {
	t.Helper()
	r := gin.New()
	r.GET("/ws/:id", func(c *gin.Context) {
		c.Set("sessionId", int64(1))
		c.Set("userId", int64(42))
		handler.HandleRoom(c)
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func wsURL(srv *httptest.Server, roomId string) string {
	return "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/" + roomId
}

func newGameState() chess.GameState {
	return *chess.NewGame(chess.NewInitialBoard())
}

// dialAndReadFirst dials, reads the first message, then returns the conn and
// the raw message bytes. The conn is NOT closed — caller decides when to close.
func dialAndReadFirst(t *testing.T, url string) (*websocket.Conn, []byte) {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)
	return conn, msg
}

// closeAndWait sends a clean close frame and waits for the server goroutine to
// finish (signalled via wg), so all deferred calls — RemoveClient etc. — have
// run before gomock's controller verifies expectations.
func closeAndWait(t *testing.T, conn *websocket.Conn, wg *sync.WaitGroup) {
	t.Helper()
	_ = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	conn.Close()
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for handler goroutine to finish")
	}
}

// makeClient returns a DoAndReturn func that builds a Client from the ws.Conn
// the handler passes in, avoiding nil-Conn panics in WriteJSON.
func makeClient(role ws.Role, color chess.PieceColor) func(string, string, ws.Conn) (*ws.Client, error) {
	return func(roomId, playerId string, conn ws.Conn) (*ws.Client, error) {
		c := color
		return &ws.Client{
			Conn:     conn,
			PlayerId: playerId,
			Role:     role,
			Color:    &c,
			Done:     make(chan struct{}),
		}, nil
	}
}

// ---------------------------------------------------------------------------
// NewHandler wiring
// ---------------------------------------------------------------------------

func TestNewHandler_DeletesRoomWhenEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	var capturedCb func(string)
	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any()).DoAndReturn(func(cb func(string)) {
		capturedCb = cb
	})
	hub.EXPECT().StartJanitor()
	roomSvc.EXPECT().DeleteRoom(gomock.Any(), "room-1").Return(nil)

	ws.NewHandler(roomSvc, userSvc, hub)

	require.NotNil(t, capturedCb)
	capturedCb("room-1")
}

func TestNewHandler_DeleteRoomError_NoPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	var capturedCb func(string)
	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any()).DoAndReturn(func(cb func(string)) {
		capturedCb = cb
	})
	hub.EXPECT().StartJanitor()
	roomSvc.EXPECT().DeleteRoom(gomock.Any(), "room-1").Return(errors.New("not found"))

	ws.NewHandler(roomSvc, userSvc, hub)

	assert.NotPanics(t, func() { capturedCb("room-1") })
}

// ---------------------------------------------------------------------------
// HandleRoom — HTTP-level rejections (no WS upgrade needed)
// ---------------------------------------------------------------------------

func TestHandleRoom_UserFetchFails_Returns500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any())
	hub.EXPECT().StartJanitor()
	userSvc.EXPECT().User(gomock.Any(), int64(42)).Return(nil, errors.New("db error"))

	handler := ws.NewHandler(roomSvc, userSvc, hub)

	r := gin.New()
	r.GET("/ws/:id", func(c *gin.Context) {
		c.Set("sessionId", int64(1))
		c.Set("userId", int64(42))
		handler.HandleRoom(c)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ws/room-1", nil))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleRoom_RoomNotFound_Returns400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any())
	hub.EXPECT().StartJanitor()
	userSvc.EXPECT().User(gomock.Any(), int64(42)).Return(&users.User{PlayerId: "p1"}, nil)
	roomSvc.EXPECT().FetchRoom(gomock.Any(), "room-1").Return(nil, errors.New("invalid room id"))

	handler := ws.NewHandler(roomSvc, userSvc, hub)

	r := gin.New()
	r.GET("/ws/:id", func(c *gin.Context) {
		c.Set("sessionId", int64(1))
		c.Set("userId", int64(42))
		handler.HandleRoom(c)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ws/room-1", nil))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// HandleRoom — full WS upgrade path
// ---------------------------------------------------------------------------

func TestHandleRoom_SpectatorConnects_ReceivesGameState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	gs := newGameState()
	var wg sync.WaitGroup
	wg.Add(1)

	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any())
	hub.EXPECT().StartJanitor()
	userSvc.EXPECT().User(gomock.Any(), int64(42)).Return(&users.User{PlayerId: "p1"}, nil)
	roomSvc.EXPECT().FetchRoom(gomock.Any(), "room-1").Return(&rooms.Room{Id: "room-1", GameState: &gs}, nil)
	hub.EXPECT().AddClient("room-1", "p1", gomock.Any()).
		DoAndReturn(makeClient(ws.Spectator, chess.PieceWhite))
	hub.EXPECT().RemoveClient("room-1", gomock.Any()).Do(func(string, ws.Conn) { wg.Done() })

	handler := ws.NewHandler(roomSvc, userSvc, hub)
	srv := newRouter(t, handler)

	conn, msg := dialAndReadFirst(t, wsURL(srv, "room-1"))

	var payload map[string]any
	require.NoError(t, json.Unmarshal(msg, &payload))
	assert.Equal(t, "game_state", payload["type"])

	closeAndWait(t, conn, &wg)
}

func TestHandleRoom_PlayerOneConnects_BroadcastsJoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	gs := newGameState()
	var wg sync.WaitGroup
	wg.Add(1)
	joinResult := contracts.PlayerJoinedResult{
		RoomId:      "room-1",
		PlayerId:    "p1",
		Success:     true,
		GameStarted: false,
	}

	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any())
	hub.EXPECT().StartJanitor()
	userSvc.EXPECT().User(gomock.Any(), int64(42)).Return(&users.User{PlayerId: "p1"}, nil)
	roomSvc.EXPECT().FetchRoom(gomock.Any(), "room-1").Return(&rooms.Room{Id: "room-1", GameState: &gs}, nil)
	hub.EXPECT().AddClient("room-1", "p1", gomock.Any()).
		DoAndReturn(makeClient(ws.PlayerOne, chess.PieceWhite))
	roomSvc.EXPECT().UpdatePlayerJoined(gomock.Any(), "room-1", "p1", chess.PieceWhite).Return(joinResult, nil)
	hub.EXPECT().Broadcast("room-1", joinResult, true).Return(nil)
	hub.EXPECT().RemoveClient("room-1", gomock.Any()).Do(func(string, ws.Conn) { wg.Done() })

	handler := ws.NewHandler(roomSvc, userSvc, hub)
	srv := newRouter(t, handler)

	conn, msg := dialAndReadFirst(t, wsURL(srv, "room-1"))

	var payload map[string]any
	require.NoError(t, json.Unmarshal(msg, &payload))
	assert.Equal(t, "game_state", payload["type"])

	closeAndWait(t, conn, &wg)
}

func TestHandleRoom_AddClientFails_ServerClosesConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomSvc := roomsMocks.NewMockIService(ctrl)
	userSvc := usersMocks.NewMockIService(ctrl)
	hub := wsMocks.NewMockIHub(ctrl)

	gs := newGameState()

	hub.EXPECT().SetOnRoomEmptyCallback(gomock.Any())
	hub.EXPECT().StartJanitor()
	userSvc.EXPECT().User(gomock.Any(), int64(42)).Return(&users.User{PlayerId: "p1"}, nil)
	roomSvc.EXPECT().FetchRoom(gomock.Any(), "room-1").Return(&rooms.Room{Id: "room-1", GameState: &gs}, nil)
	hub.EXPECT().AddClient("room-1", "p1", gomock.Any()).Return(nil, errors.New("room full"))

	handler := ws.NewHandler(roomSvc, userSvc, hub)
	srv := newRouter(t, handler)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL(srv, "room-1"), nil)
	if err == nil {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _, readErr := conn.ReadMessage()
		assert.Error(t, readErr)
		conn.Close()
	}
}
