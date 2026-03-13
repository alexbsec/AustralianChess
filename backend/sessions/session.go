package sessions

import "context"

type Repository interface {
	NewSession(ctx context.Context, userId int64, token string) (*Session, error)
	GetUserSession(ctx context.Context, userId int64) (*Session, error)
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	UpdateSession(ctx context.Context, session *Session) (*Session, error)
	DestroySession(ctx context.Context, id int64) error
}

type IService interface {
	CreateSession(ctx context.Context, userId int64) (*Session, error)
	ValidateSession(ctx context.Context, token string) (*Session, bool, error)
	RevokeSession(ctx context.Context, token string) error
}
