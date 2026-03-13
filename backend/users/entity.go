package users

import "time"

type User struct {
	Id           int64     `db:"id"`
	PlayerId     string    `db:"player_id"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type MinifiedUser struct {
	Id       int64  `json:"id"`
	PlayerId string `json:"player_id"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        MinifiedUser `json:"user"`
}

type NewUserResponse struct {
	Message  string `json:"message"`
	PlayerId string `json:"player_id"`
}
