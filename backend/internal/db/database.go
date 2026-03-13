package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Connect(ctx context.Context) error
	Pool() *pgxpool.Pool
	Close()
}
