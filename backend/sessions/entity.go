package sessions

import "time"

type Session struct {
	Id           int64      `db:"id"`
	UserId       int64      `db:"user_id"`
	SessionToken string     `db:"session_token"`
	ExpiresAt    time.Time  `db:"expires_at"`
	RevokedAt    *time.Time `db:"revoked_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

