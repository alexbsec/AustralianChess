package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/db/repositories"
	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/stretchr/testify/require"
)

func TestSessionRepository_CRUD(t *testing.T) {
	db := repositories.SetupDB(t)
	defer db.Close()
	repositories.CleanDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	ctx := context.Background()

	user := &users.User{
		PlayerId:     "player-session",
		PasswordHash: "hash",
	}

	createdUser, err := userRepo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, createdUser.Id)

	userId := createdUser.Id

	session, err := sessionRepo.NewSession(ctx, userId, "token123")
	require.NoError(t, err)
	require.NotZero(t, session.Id)
	require.Equal(t, userId, session.UserId)

	userSession, err := sessionRepo.GetUserSession(ctx, userId)
	require.NoError(t, err)
	require.Equal(t, session.Id, userSession.Id)

	tokenSession, err := sessionRepo.GetSessionByToken(ctx, "token123")
	require.NoError(t, err)
	require.Equal(t, session.Id, tokenSession.Id)

	oldUpdatedAt := session.UpdatedAt
	time.Sleep(1 * time.Second)

	session.SessionToken = "token456"
	session.ExpiresAt = time.Now().Add(48 * time.Hour)

	updatedSession, err := sessionRepo.UpdateSession(ctx, session)
	require.NoError(t, err)
	require.Equal(t, "token456", updatedSession.SessionToken)
	require.True(t, updatedSession.UpdatedAt.After(oldUpdatedAt))

	err = sessionRepo.DestroySession(ctx, session.Id)
	require.NoError(t, err)

	_, err = sessionRepo.GetUserSession(ctx, userId)
	require.Error(t, err)

	_, err = sessionRepo.GetSessionByToken(ctx, "token456")
	require.Error(t, err)
}
