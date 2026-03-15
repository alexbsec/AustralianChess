package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/alexbsec/AustralianChess/backend/sessions"
)

type Authorizer struct {
	sessionService sessions.IService
}

func NewAuthorizer(sessionService sessions.IService) *Authorizer {
	return &Authorizer{
		sessionService: sessionService,
	}
}

func (a *Authorizer) Authorize(ctx context.Context, request *http.Request) (*sessions.Session, error) {
	token, err := extractToken(request)
	if err != nil {
		return nil, err
	}

	session, valid, err := a.sessionService.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, ErrInvalidToken
	}

	return session, nil
}

func isBearerToken(token string) bool {
	return len(token) > 7 && token[:7] == "Bearer "
}

func extractTokenBearer(token string) string {
	return token[7:]
}

func extractToken(request *http.Request) (string, error) {
	token := request.Header.Get("Authorization")
	if token != "" {
		if !isBearerToken(token) {
			log.Printf("invalid token format: %s. Expected token Bearer", token)
			return "", ErrInvalidToken
		}

		return extractTokenBearer(token), nil
	}

	// query param fallback
	token = request.URL.Query().Get("token")
	if token != "" {
		return token, nil
	}

	return "", ErrMissingToken
}
