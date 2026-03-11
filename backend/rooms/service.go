package rooms

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/google/uuid"
)

type Service struct {
	mtx   sync.RWMutex
	rooms map[string]*Room
}

func NewService() *Service {
	return &Service{
		rooms: make(map[string]*Room),
	}
}

func (s *Service) CreateRoom(ctx context.Context) (CreateRoomResponse, error) {
	res := CreateRoomResponse{}

	board := chess.NewInitialBoard()
	gameState := chess.NewGame(board)

	room := &Room{
		Id:        uuid.New().String(),
		GameState: gameState,
		CreatedAt: time.Now(),
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.rooms[room.Id] = room

	res.RoomId = room.Id
	return res, nil
}

func (s *Service) DisplayRoom(ctx context.Context, roomId string) (DisplayRoomResponse, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	room, ok := s.rooms[roomId]
	if !ok {
		return DisplayRoomResponse{}, errors.New("invalid room")
	}

	res := DisplayRoomResponse{
		RoomId:    roomId,
		Status:    "waiting",
		GameState: *room.GameState,
	}

	return res, nil
}
