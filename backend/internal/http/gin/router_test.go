package gin_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ginAuss "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	authMocks "github.com/alexbsec/AustralianChess/backend/internal/auth/mocks"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	roomsMocks "github.com/alexbsec/AustralianChess/backend/rooms/mocks"
	usersMocks "github.com/alexbsec/AustralianChess/backend/users/mocks"
	"github.com/alexbsec/AustralianChess/backend/sessions"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)



func TestMakeHandlers_FullSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomService := roomsMocks.NewMockIService(ctrl)
	userService := usersMocks.NewMockIService(ctrl)
	authService := authMocks.NewMockIAuthorizer(ctrl)

	// auth middleware will always pass
	authService.EXPECT().
		Authorize(gomock.Any(), gomock.Any()).
		Return(&sessions.Session{Id: 2, UserId: 1}, nil).
		AnyTimes()

	roomService.EXPECT().
		CreateRoom(gomock.Any()).
		Return(rooms.CreateRoomResponse{}, errors.New("fail"))

	router := ginAuss.MakeHandlers(
		context.Background(),
		roomService,
		userService,
		nil, // wsHandler
		authService,
	)

	// --- test that routes exist ---
	req := httptest.NewRequest(http.MethodGet, "/api/v1/room/create", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// since roomService returns error → 500, but route exists
	require.NotEqual(t, http.StatusNotFound, resp.Code)
}


func TestMakeHandlers_UserRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := usersMocks.NewMockIService(ctrl)
	authService := authMocks.NewMockIAuthorizer(ctrl)

	userService.EXPECT().
		LoginUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("fail"))

	router := ginAuss.MakeHandlers(
		context.Background(),
		nil,
		userService,
		nil,
		authService,
	)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBufferString(`{}`))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.NotEqual(t, http.StatusNotFound, resp.Code)
}


func TestMakeHandlers_RoomRequiresAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomService := roomsMocks.NewMockIService(ctrl)
	authService := authMocks.NewMockIAuthorizer(ctrl)

	// simulate auth failure
	authService.EXPECT().
		Authorize(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("unauthorized"))

	router := ginAuss.MakeHandlers(
		context.Background(),
		roomService,
		nil,
		nil,
		authService,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/room/create", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}


func TestMakeHandlers_RoomRouteExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomService := roomsMocks.NewMockIService(ctrl)
	authService := authMocks.NewMockIAuthorizer(ctrl)

	authService.EXPECT().
		Authorize(gomock.Any(), gomock.Any()).
		Return(&sessions.Session{Id: 2, UserId: 1}, nil).
		AnyTimes()

	roomService.EXPECT().
		DisplayRoom(gomock.Any(), "123").
		Return(rooms.DisplayRoomResponse{}, errors.New("not found"))

	router := ginAuss.MakeHandlers(
		context.Background(),
		roomService,
		nil,
		nil,
		authService,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/room/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.NotEqual(t, http.StatusNotFound, resp.Code)
}


func TestMakeHandlers_NilServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := authMocks.NewMockIAuthorizer(ctrl)

	router := ginAuss.MakeHandlers(
		context.Background(),
		nil,
		nil,
		nil,
		authService,
	)

	// should still create router without panic
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.NotEqual(t, 0, resp.Code)
}
