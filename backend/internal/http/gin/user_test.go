package gin_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	ginAuss "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/alexbsec/AustralianChess/backend/users/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMakeUserLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	mockService.EXPECT().
		LoginUser(gomock.Any(), "user", "pass").
		Return(&users.LoginResponse{
			AccessToken: "jwt-token",
			User: users.MinifiedUser{
				Id:       1,
				PlayerId: "player-123",
			},
		}, nil)

	r := gin.New()
	r.POST("/login", ginAuss.MakeUserLogin(mockService))

	body := `{"username":"user","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	require.Contains(t, resp.Body.String(), "jwt-token")
	require.Contains(t, resp.Body.String(), "player-123")
}


func TestMakeUserLogin_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	r := gin.New()
	r.POST("/login", ginAuss.MakeUserLogin(mockService))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`invalid json`))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}


func TestMakeUserLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	mockService.EXPECT().
		LoginUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, users.ErrInvalidCredentials)

	r := gin.New()
	r.POST("/login", ginAuss.MakeUserLogin(mockService))

	body := `{"username":"user","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}


func TestMakeUserRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	mockService.EXPECT().
		NewUser(gomock.Any(), "user", "pass").
		Return(&users.NewUserResponse{
			Message:  "user created",
			PlayerId: "player-123",
		}, nil)

	r := gin.New()
	r.POST("/register", ginAuss.MakeUserRegister(mockService))

	body := `{
		"username":"user",
		"password":"pass",
		"confirmPassword":"pass"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	require.Contains(t, resp.Body.String(), "user created")
	require.Contains(t, resp.Body.String(), "player-123")
}


func TestMakeUserRegister_PasswordMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	r := gin.New()
	r.POST("/register", ginAuss.MakeUserRegister(mockService))

	body := `{
		"username":"user",
		"password":"pass",
		"confirmPassword":"different"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "passwords do not match")
}


func TestMakeUserRegister_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIService(ctrl)

	mockService.EXPECT().
		NewUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, users.ErrPasswordTooShort)

	r := gin.New()
	r.POST("/register", ginAuss.MakeUserRegister(mockService))

	body := `{
		"username":"user",
		"password":"123",
		"confirmPassword":"123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}
