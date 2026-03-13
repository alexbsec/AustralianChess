package db

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	connString string
	pool       *pgxpool.Pool
}

func NewPostgresDB(connString string) *PostgresDB {
	return &PostgresDB{
		connString: connString,
	}
}

func (db *PostgresDB) Connect(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, db.connString)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	db.pool = pool
	return nil
}

func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func PostgresConnStringFromEnv(envVars map[string]any) string {
	host := envVars["DB_HOST"].(string)
	port := strconv.Itoa(envVars["DB_PORT"].(int))
	username := envVars["DB_USERNAME"].(string)
	password := envVars["DB_PASSWORD"].(string)
	dbName := envVars["DB_NAME"].(string)

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbName,
	}

	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	return u.String()
}
