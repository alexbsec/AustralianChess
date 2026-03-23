package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/db/repositories"
	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateUser(t *testing.T) {
	database := repositories.SetupDB(t)
	defer database.Close()
	repositories.CleanDB(t, database)

	repo := repositories.NewUserRepository(database)

	user := &users.User{
		PlayerId:     "p1",
		PasswordHash: "hash",
	}

	created, err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotZero(t, created.Id)
}

func TestUserRepository_CRUD(t *testing.T) {
	db := repositories.SetupDB(t)
	defer db.Close()
	repositories.CleanDB(t, db)

	repo := repositories.NewUserRepository(db)

	ctx := context.Background()

	user := &users.User{
		PlayerId:     "player1",
		PasswordHash: "hash123",
	}
	created, err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, created.Id)
	require.Equal(t, "player1", created.PlayerId)

	fetched, err := repo.FetchUserByPlayerId(ctx, "player1")
	require.NoError(t, err)
	require.Equal(t, created.Id, fetched.Id)
	require.Equal(t, created.PasswordHash, fetched.PasswordHash)

	fetchedByID, err := repo.FetchUser(ctx, created.Id)
	require.NoError(t, err)
	require.Equal(t, created.PlayerId, fetchedByID.PlayerId)

	oldUpdatedAt := created.UpdatedAt
	time.Sleep(1 * time.Second) // ensure timestamp change
	created.PasswordHash = "newhash"
	updated, err := repo.UpdateUser(ctx, created)
	require.NoError(t, err)
	require.Equal(t, "newhash", updated.PasswordHash)
	require.True(t, updated.UpdatedAt.After(oldUpdatedAt))

	err = repo.DeleteUser(ctx, created.Id)
	require.NoError(t, err)

	_, err = repo.FetchUser(ctx, created.Id)
	require.Error(t, err)
}
