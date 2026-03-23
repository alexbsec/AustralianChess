package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/db"
	"github.com/stretchr/testify/require"
)


func TestPostgresConnStringFromEnv(t *testing.T) {
	env := map[string]any{
		"DB_HOST":     "localhost",
		"DB_PORT":     5432,
		"DB_USERNAME": "user",
		"DB_PASSWORD": "pass",
		"DB_NAME":     "testdb",
	}

	conn := db.PostgresConnStringFromEnv(env)

	require.Contains(t, conn, "postgres://user:pass@localhost:5432/testdb")
	require.Contains(t, conn, "sslmode=disable")
}


func TestPostgresConnStringFromEnv_Exact(t *testing.T) {
	env := map[string]any{
		"DB_HOST":     "localhost",
		"DB_PORT":     5432,
		"DB_USERNAME": "user",
		"DB_PASSWORD": "pass",
		"DB_NAME":     "testdb",
	}

	conn := db.PostgresConnStringFromEnv(env)

	require.Equal(
		t,
		"postgres://user:pass@localhost:5432/testdb?sslmode=disable",
		conn,
	)
}


func TestPostgresDB_Connect(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")

	db := db.NewPostgresDB(connStr)

	err := db.Connect(context.Background())
	require.NoError(t, err)

	require.NotNil(t, db.Pool())

	db.Close()
}


func TestPostgresDB_Close(t *testing.T) {
	connStr := os.Getenv("DATABASE_URL")

	db := db.NewPostgresDB(connStr)

	err := db.Connect(context.Background())
	require.NoError(t, err)

	db.Close()

	// after close, pool should be nil (or at least closed safely)
	require.NotPanics(t, func() {
		db.Close()
	})
}


func TestPostgresDB_Connect_Invalid(t *testing.T) {
	db := db.NewPostgresDB("invalid-connection-string")

	err := db.Connect(context.Background())
	require.Error(t, err)
}
