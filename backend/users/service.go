package users

import (
	"context"
	"strings"

	"github.com/alexbsec/AustralianChess/backend/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
	sessionService sessions.IService
}

func NewService(repo Repository, sessionService sessions.IService) *Service {
	return &Service{
		repo: repo,
		sessionService: sessionService,
	}
}

func (s *Service) LoginUser(ctx context.Context, playerId, password string) (*LoginResponse, error) {
	playerId, err := sanctifyPlayerId(playerId)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FetchUserByPlayerId(ctx, playerId)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	session, err := s.sessionService.CreateSession(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		AccessToken: session.SessionToken,
		User: MinifiedUser{
			Id: user.Id,
			PlayerId: user.PlayerId,
		},
	}

	return response, nil
}

func (s *Service) NewUser(ctx context.Context, playerId, password string) (*NewUserResponse, error) {
	playerId, err := sanctifyPlayerId(playerId)
	if err != nil {
		return nil, err
	}

	currUser, err := s.repo.FetchUserByPlayerId(ctx, playerId)
	if err == nil && currUser != nil {
		return nil, ErrPlayerIdAlreadyExists
	}

	err = validatePassword(password)
	if err != nil {
		return nil, err
	}

	hashPwd, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		PlayerId: playerId,
		PasswordHash: hashPwd,
	}

	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	response := &NewUserResponse{
		Message: "account created successfully",
		PlayerId: playerId,
	}

	return response, nil
}

func hashPassword(password string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPwd), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func sanctifyPlayerId(playerId string) (string, error) {
	playerId = strings.TrimSpace(playerId)
	for _, r := range playerId {
		if !(r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return "", ErrInvalidPlayerId
		}
	}

	return playerId, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
