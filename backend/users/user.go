package users

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FetchUserByPlayerId(ctx context.Context, playerId string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, userId int64) error
}

type IService interface {
	LoginUser(ctx context.Context, playerId, password string) (*LoginResponse, error)
	NewUser(ctx context.Context, playerId, password string) (*NewUserResponse, error)
}

