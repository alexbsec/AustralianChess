package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSession(ctx context.Context, userId int64) (*Session, error) {
	sessionToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	session, err := s.repo.NewSession(ctx, userId, sessionToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) ValidateSession(ctx context.Context, token string) (*Session, bool, error) {
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, false, err
	}

	now := time.Now().UTC()
	if session.ExpiresAt.Before(now) || session.RevokedAt != nil {
		return session, false, nil
	}

	return session, true, nil
}

func (s *Service) RevokeSession(ctx context.Context, token string) error {
	session, valid, err := s.ValidateSession(ctx, token)
	if err != nil {
		return err
	}

	if !valid {
		return nil
	}

	return s.repo.DestroySession(ctx, session.Id)
}

func (s *Service) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
