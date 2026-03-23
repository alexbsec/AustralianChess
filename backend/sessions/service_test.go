package sessions_test

import (
	"context"
	"errors"
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

func TestValidateSession_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "nonexistent").
		Return(nil, errors.New("not found"))

	svc := sessions.NewService(mockRepo)
	_, valid, err := svc.ValidateSession(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent token")
	}
	if valid {
		t.Fatal("expected invalid session")
	}
}

func TestRevokeSession_AlreadyExpired(t *testing.T) {
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

	// DestroySession should NOT be called for an already invalid session
	svc := sessions.NewService(mockRepo)
	err := svc.RevokeSession(context.Background(), "token")
	if err != nil {
		t.Fatalf("unexpected error revoking expired session: %v", err)
	}
}

func TestRevokeSession_AlreadyRevoked(t *testing.T) {
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

	// DestroySession should NOT be called
	svc := sessions.NewService(mockRepo)
	err := svc.RevokeSession(context.Background(), "token")
	if err != nil {
		t.Fatalf("unexpected error revoking already-revoked session: %v", err)
	}
}

func TestRevokeSession_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "token").
		Return(nil, errors.New("db error"))

	svc := sessions.NewService(mockRepo)
	err := svc.RevokeSession(context.Background(), "token")
	if err == nil {
		t.Fatal("expected error from repo failure")
	}
}

func TestRevokeSession_DestroyFails(t *testing.T) {
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
		Return(errors.New("db error"))

	svc := sessions.NewService(mockRepo)
	err := svc.RevokeSession(context.Background(), "token")
	if err == nil {
		t.Fatal("expected error when destroy fails")
	}
}

func TestCreateSession_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.
		EXPECT().
		NewSession(gomock.Any(), int64(123), gomock.Any()).
		Return(nil, errors.New("db error"))

	svc := sessions.NewService(mockRepo)
	session, err := svc.CreateSession(context.Background(), 123)
	if err == nil {
		t.Fatal("expected error from repo failure")
	}
	if session != nil {
		t.Fatal("expected nil session on error")
	}
}

func TestValidateSession_ExpiresExactlyNow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.
		EXPECT().
		GetSessionByToken(gomock.Any(), "token").
		Return(&sessions.Session{
			Id:           1,
			SessionToken: "token",
			ExpiresAt:    time.Now().Add(-time.Millisecond), // just expired
		}, nil)

	svc := sessions.NewService(mockRepo)
	_, valid, err := svc.ValidateSession(context.Background(), "token")
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid for just-expired session")
	}
}

func TestCreateSession_TokenIsUnique(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	tokens := make(map[string]bool)

	mockRepo.
		EXPECT().
		NewSession(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, userId int64, token string) (*sessions.Session, error) {
			tokens[token] = true
			return &sessions.Session{
				Id:           1,
				UserId:       userId,
				SessionToken: token,
				ExpiresAt:    time.Now().Add(time.Hour),
			}, nil
		}).
		Times(10)

	svc := sessions.NewService(mockRepo)
	for i := range 10 {
		_, err := svc.CreateSession(context.Background(), int64(i))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if len(tokens) != 10 {
		t.Fatalf("expected 10 unique tokens, got %d", len(tokens))
	}
}
