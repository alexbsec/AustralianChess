package auth

import (
	"context"
	"net/http"

	"github.com/alexbsec/AustralianChess/backend/sessions"
)

type IAuthorizer interface {
	Authorize(ctx context.Context, request *http.Request) (*sessions.Session, error)
}
