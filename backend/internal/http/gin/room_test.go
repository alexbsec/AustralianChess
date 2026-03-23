package gin_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ginAuss "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
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
