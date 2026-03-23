package sessions_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/sessions"
	"github.com/alexbsec/AustralianChess/backend/sessions/mocks"
	"github.com/golang/mock/gomock"
)

func TestCreateSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.
		EXPECT().
		NewSession(gomock.Any(), int64(123), gomock.Any()).
		DoAndReturn(func(ctx context.Context, userId int64, token string) (*sessions.Session, error) {
			if token == "" {
				t.Fatal("expected token to be generated")
			}

			return &sessions.Session{
				Id:           1,
				UserId:       userId,
				SessionToken: token,
				ExpiresAt:    time.Now().Add(time.Hour),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		})

	svc := sessions.NewService(mockRepo)

	session, err := svc.CreateSession(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session == nil {
		t.Fatal("expected session")
	}

	if session.SessionToken == "" {
		t.Fatal("expected session token to be set")
	}
}


func TestValidateSession_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "token").
		Return(&sessions.Session{
			Id:           1,
			UserId:       123,
			SessionToken: "token",
			ExpiresAt:    time.Now().Add(time.Hour),
			RevokedAt:    nil,
		}, nil)

	svc := sessions.NewService(mockRepo)

	session, valid, err := svc.ValidateSession(context.Background(), "token")

	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("expected valid session")
	}

	if session == nil || session.SessionToken != "token" {
		t.Fatal("unexpected session returned")
	}
}


func TestValidateSession_Expired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "token").
		Return(&sessions.Session{
			Id:           1,
			SessionToken: "token",
			ExpiresAt:    time.Now().Add(-time.Hour),
		}, nil)

	svc := sessions.NewService(mockRepo)

	_, valid, err := svc.ValidateSession(context.Background(), "token")

	if err != nil {
		t.Fatal(err)
	}

	if valid {
		t.Fatal("expected invalid (expired)")
	}
}


func TestValidateSession_Revoked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	now := time.Now()

	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "token").
		Return(&sessions.Session{
			Id:           1,
			SessionToken: "token",
			ExpiresAt:    time.Now().Add(time.Hour),
			RevokedAt:    &now,
		}, nil)

	svc := sessions.NewService(mockRepo)

	_, valid, err := svc.ValidateSession(context.Background(), "token")

	if err != nil {
		t.Fatal(err)
	}

	if valid {
		t.Fatal("expected invalid (revoked)")
	}
}


func TestRevokeSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "token").
		Return(&sessions.Session{
			Id:           1,
			SessionToken: "token",
			ExpiresAt:    time.Now().Add(time.Hour),
		}, nil)

	mockRepo.
		EXPECT().
		DestroySession(gomock.Any(), int64(1)).
		Return(nil)

	svc := sessions.NewService(mockRepo)

	err := svc.RevokeSession(context.Background(), "token")

	if err != nil {
		t.Fatal(err)
	}
}
