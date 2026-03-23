package repositories

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/db"
	"github.com/stretchr/testify/require"
)

func SetupDB(t *testing.T) *db.PostgresDB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")

	// Run migrations
	cmd := exec.Command("goose", "-dir", "../../../migrations", "postgres", dsn, "up")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	database := db.NewPostgresDB(dsn)
	err = database.Connect(context.Background())
	require.NoError(t, err)

	return database
}

func CleanDB(t *testing.T, database *db.PostgresDB) {
	t.Helper()
	_, err := database.Pool().Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}
