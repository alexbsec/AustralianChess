package repositories

import (
	"context"

	"github.com/alexbsec/AustralianChess/backend/internal/db"
	"github.com/alexbsec/AustralianChess/backend/users"
)

type UserRepository struct {
	database db.Database
}

func NewUserRepository(database db.Database) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *users.User) (*users.User, error) {
	query := `
	INSERT INTO users (player_id, password_hash)
	VALUES ($1, $2)
	RETURNING id, created_at, updated_at
	`
	pool := r.database.Pool()	
	row := pool.QueryRow(ctx, query, user.PlayerId, user.PasswordHash)
	err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *UserRepository) FetchUserByPlayerId(ctx context.Context, playerId string) (*users.User, error) {
	query := `
	SELECT id, player_id, password_hash, created_at, updated_at
	FROM users
	WHERE player_id = $1
	`

	pool := r.database.Pool()	
	row := pool.QueryRow(ctx, query, playerId)
	var user users.User
	err := row.Scan(&user.Id, &user.PlayerId, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *users.User) (*users.User, error) {
	query := `
	UPDATE users
	SET password_hash = $1, updated_at = NOW()
	WHERE id = $2
	RETURNING updated_at
	`

	pool := r.database.Pool()	
	row := pool.QueryRow(ctx, query, user.PasswordHash, user.Id)
	err := row.Scan(&user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId int64) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	pool := r.database.Pool()
	_, err := pool.Exec(ctx, query, userId)
	return err
}
