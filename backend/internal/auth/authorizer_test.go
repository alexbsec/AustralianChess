package auth_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/auth"
	"github.com/alexbsec/AustralianChess/backend/sessions"
	sessionMock "github.com/alexbsec/AustralianChess/backend/sessions/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAuthorize_Success_BearerToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := sessionMock.NewMockIService(ctrl)

	auth := auth.NewAuthorizer(mockSvc)

	req := &http.Request{
		Header: make(http.Header),
	}

	req.Header.Set("Authorization", "Bearer valid-token")

	session := &sessions.Session{Id: 1}

	mockSvc.
		EXPECT().
		ValidateSession(gomock.Any(), "valid-token").
		Return(session, true, nil)

	result, err := auth.Authorize(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, session, result)
}

func TestAuthorize_InvalidTokenFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := sessionMock.NewMockIService(ctrl)

	auth := auth.NewAuthorizer(mockSvc)

	req := &http.Request{
		Header: make(http.Header),
	}

	req.Header.Set("Authorization", "InvalidToken")

	result, err := auth.Authorize(context.Background(), req)

	require.Nil(t, result)
	require.Error(t, err)
}


func TestAuthorize_MissingToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := sessionMock.NewMockIService(ctrl)

	auth := auth.NewAuthorizer(mockSvc)

	req := &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	result, err := auth.Authorize(context.Background(), req)

	require.Nil(t, result)
	require.Error(t, err)
}
