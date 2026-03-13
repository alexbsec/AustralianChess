package repositories

import (
	"context"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/db"
	"github.com/alexbsec/AustralianChess/backend/sessions"
)

type SessionRepository struct {
	database db.Database
}

func NewSessionRepository(database db.Database) *SessionRepository {
	return &SessionRepository{
		database: database,
	}
}

func (r *SessionRepository) NewSession(
	ctx context.Context,
	userId int64,
	token string,
) (*sessions.Session, error) {
	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	query := `
	INSERT INTO sessions (user_id, token, expires_at)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, token, expires_at, created_at, updated_at
	`

	pool := r.database.Pool()
	var session sessions.Session
	err := pool.QueryRow(ctx, query, userId, token, expiresAt).
		Scan(
			&session.Id,
			&session.UserId,
			&session.SessionToken,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetUserSession(ctx context.Context, userId int64) (*sessions.Session, error) {
	query := `
	SELECT id, user_id, token, expires_at, revoked_at, created_at, updated_at
	FROM sessions
	WHERE user_id = $1
	  AND revoked_at IS NULL
	  AND expires_at > NOW()
	ORDER BY created_at DESC
	LIMIT 1
	`

	pool := r.database.Pool()
	var session sessions.Session
	err := pool.QueryRow(ctx, query, userId).
		Scan(
			&session.Id,
			&session.UserId,
			&session.SessionToken,
			&session.ExpiresAt,
			&session.RevokedAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetSessionByToken(ctx context.Context, token string) (*sessions.Session, error) {
	query := `
	SELECT id, user_id, token, expires_at, revoked_at, created_at, updated_at
	FROM sessions
	WHERE token = $1
	  AND revoked_at IS NULL
	  AND expires_at > NOW()
	`

	pool := r.database.Pool()
	var session sessions.Session
	err := pool.QueryRow(ctx, query, token).Scan(
		&session.Id,
		&session.UserId,
		&session.SessionToken,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *sessions.Session) (*sessions.Session, error) {
	updatedAt := time.Now().UTC()
	query := `
	UPDATE sessions
	SET token = $1, expires_at = $2, revoked_at = $3, updated_at = $4
	WHERE id = $5
	RETURNING updated_at
	`

	pool := r.database.Pool()
	err := pool.QueryRow(ctx, query,
		session.SessionToken,
		session.ExpiresAt,
		session.RevokedAt,
		updatedAt,
		session.Id,
	).Scan(&session.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) DestroySession(ctx context.Context, id int64) error {
	revokedAt := time.Now().UTC()

	query := `
	UPDATE sessions
	SET revoked_at = $1, updated_at = $2
	WHERE id = $3
	`

	pool := r.database.Pool()
	_, err := pool.Exec(ctx, query, revokedAt, revokedAt, id)
	return err
}
