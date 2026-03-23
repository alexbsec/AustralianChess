package gin_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/auth/mocks"
	ginAuss "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	"github.com/alexbsec/AustralianChess/backend/sessions"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupAuthTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware)

	r.GET("/test", func(ctx *gin.Context) {
		userId, _ := ctx.Get("userId")
		sessionId, _ := ctx.Get("sessionId")

		ctx.JSON(http.StatusOK, gin.H{
			"userId":    userId,
			"sessionId": sessionId,
		})
	})

	return r
}

func TestAuthMiddleware_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockIAuthorizer(ctrl)

	session := &sessions.Session{
		Id:        456,
		UserId:    123,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockAuth.EXPECT().
		Authorize(gomock.Any(), gomock.Any()).
		Return(session, nil)

	router := setupAuthTestRouter(ginAuss.AuthMiddleware(mockAuth))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)

	require.Equal(t, float64(123), body["userId"])    
	require.Equal(t, float64(456), body["sessionId"]) 
}

func TestAuthMiddleware_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockIAuthorizer(ctrl)

	mockAuth.EXPECT().
		Authorize(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("unauthorized"))

	router := setupAuthTestRouter(ginAuss.AuthMiddleware(mockAuth))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}
