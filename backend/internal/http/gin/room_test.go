package gin_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	ginAuss "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	wsMocks "github.com/alexbsec/AustralianChess/backend/internal/ws/mocks"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/alexbsec/AustralianChess/backend/rooms/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMakeNewRoom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	mockService.EXPECT().
		CreateRoom(gomock.Any()).
		Return(
			rooms.CreateRoomResponse{RoomId: "room-123"},
			nil,
		)

	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("userId", int64(1))
		ctx.Set("sessionId", int64(2))
		ctx.Next()
	})

	r.GET("/create", ginAuss.MakeNewRoom(mockService))

	req := httptest.NewRequest(http.MethodGet, "/create", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "room-123")
}

func TestMakeNewRoom_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	mockService.EXPECT().
		CreateRoom(gomock.Any()).
		Return(rooms.CreateRoomResponse{}, errors.New("fail"))

	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("userId", int64(1))
		ctx.Set("sessionId", int64(2))
		ctx.Next()
	})

	r.GET("/create", ginAuss.MakeNewRoom(mockService))

	req := httptest.NewRequest(http.MethodGet, "/create", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestMultiplayer_MissingRoomId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := wsMocks.NewMockIHandler(ctrl)

	r := gin.New()
	r.GET("/room/:id", ginAuss.Multiplayer(mockHandler))

	req := httptest.NewRequest(http.MethodGet, "/room/", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}


func TestMultiplayer_EmptyRoomId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := wsMocks.NewMockIHandler(ctrl)

	r := gin.New()
	r.GET("/room/:id", func(ctx *gin.Context) {
		ctx.Params = gin.Params{{Key: "id", Value: ""}}
		ginAuss.Multiplayer(mockHandler)(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, "/room/anything", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "missing room id")
}


func TestMultiplayer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := wsMocks.NewMockIHandler(ctrl)

	mockHandler.EXPECT().
		HandleRoom(gomock.Any(), "room-123").
		Times(1)

	r := gin.New()
	r.GET("/room/:id", ginAuss.Multiplayer(mockHandler))

	req := httptest.NewRequest(http.MethodGet, "/room/room-123", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}


func TestSingleplayer_MissingRoomId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := wsMocks.NewMockIHandler(ctrl)

	r := gin.New()
	r.GET("/single", ginAuss.Singleplayer(mockHandler))

	req := httptest.NewRequest(http.MethodGet, "/single", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "missing roomId")
}


func TestSingleplayer_InvalidDifficulty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := wsMocks.NewMockIHandler(ctrl)

	r := gin.New()
	r.GET("/single", ginAuss.Singleplayer(mockHandler))

	req := httptest.NewRequest(http.MethodGet, "/single?roomId=room-1&difficulty=abc", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "invalid difficulty")
}


func TestSingleplayer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := wsMocks.NewMockIHandler(ctrl)

	mockHandler.EXPECT().
		PlayBot(gomock.Any(), gomock.AssignableToTypeOf(ws.PlayBotDTO{})).
		Do(func(_ *gin.Context, dto ws.PlayBotDTO) {
			require.Equal(t, "room-1", dto.RoomId)
			require.Equal(t, bot.Difficulty(2), dto.Difficulty)
			require.Equal(t, chess.PieceColor(1), dto.PlayerPlayingAs)
		}).
		Times(1)

	r := gin.New()
	r.GET("/single", ginAuss.Singleplayer(mockHandler))

	req := httptest.NewRequest(
		http.MethodGet,
		"/single?roomId=room-1&difficulty=2&player_playing_as=1",
		nil,
	)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}
