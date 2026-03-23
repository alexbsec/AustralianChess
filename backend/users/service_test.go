package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/sessions"
	sessionMocks "github.com/alexbsec/AustralianChess/backend/sessions/mocks"
	"github.com/alexbsec/AustralianChess/backend/users"
	userMocks "github.com/alexbsec/AustralianChess/backend/users/mocks"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(&users.User{
			Id:           1,
			PlayerId:     "player1",
			PasswordHash: string(hashedPwd),
		}, nil)

	mockSession.
		EXPECT().
		CreateSession(gomock.Any(), int64(1)).
		Return(&sessions.Session{
			SessionToken: "token123",
		}, nil)

	svc := users.NewService(mockRepo, mockSession)

	resp, err := svc.LoginUser(context.Background(), "player1", "password123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AccessToken != "token123" {
		t.Fatal("expected correct token")
	}

	if resp.User.PlayerId != "player1" {
		t.Fatal("wrong user returned")
	}
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(&users.User{
			Id:           1,
			PlayerId:     "player1",
			PasswordHash: string(hashedPwd),
		}, nil)

	svc := users.NewService(mockRepo, mockSession)

	_, err := svc.LoginUser(context.Background(), "player1", "wrong-password")

	if err != users.ErrInvalidCredentials {
		t.Fatal("expected invalid credentials error")
	}
}

func TestLoginUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(nil, errors.New("db error"))

	svc := users.NewService(mockRepo, mockSession)

	_, err := svc.LoginUser(context.Background(), "player1", "password")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoginUser_SessionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(&users.User{
			Id:           1,
			PlayerId:     "player1",
			PasswordHash: string(hashedPwd),
		}, nil)

	mockSession.
		EXPECT().
		CreateSession(gomock.Any(), int64(1)).
		Return(nil, errors.New("session error"))

	svc := users.NewService(mockRepo, mockSession)

	_, err := svc.LoginUser(context.Background(), "player1", "password123")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(nil, errors.New("not found")) // simulate user not existing

	mockRepo.
		EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *users.User) (*users.User, error) {
			if user.PasswordHash == "" {
				t.Fatal("expected hashed password")
			}
			return user, nil
		})

	svc := users.NewService(mockRepo, mockSession)

	resp, err := svc.NewUser(context.Background(), "player1", "password123")

	if err != nil {
		t.Fatal(err)
	}

	if resp.PlayerId != "player1" {
		t.Fatal("wrong playerId")
	}
}

func TestNewUser_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(&users.User{Id: 1}, nil)

	svc := users.NewService(mockRepo, mockSession)

	_, err := svc.NewUser(context.Background(), "player1", "password123")

	if err != users.ErrPlayerIdAlreadyExists {
		t.Fatal("expected already exists error")
	}
}

func TestNewUser_WeakPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMocks.NewMockRepository(ctrl)
	mockSession := sessionMocks.NewMockIService(ctrl)

	mockRepo.
		EXPECT().
		FetchUserByPlayerId(gomock.Any(), "player1").
		Return(nil, errors.New("not found"))

	svc := users.NewService(mockRepo, mockSession)

	_, err := svc.NewUser(context.Background(), "player1", "123")

	if err != users.ErrPasswordTooShort {
		t.Fatal("expected password too short error")
	}
}
