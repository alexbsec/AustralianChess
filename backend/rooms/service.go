package rooms

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/ws/parser"
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

func (s *Service) FetchRoom(ctx context.Context, roomId string) (*Room, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	room, ok := s.rooms[roomId]
	if !ok {
		return nil, errors.New("invalid room id")
	}

	return room, nil
}

func (s *Service) ExecuteCommand(ctx context.Context, cmd parser.Command) (parser.Result, error) {
	if cmd == nil {
		return nil, errors.New("cannot execute nil command")
	}

	switch cmd.(type) {
	case parser.MoveCommand:
		moveCmd := cmd.(parser.MoveCommand)
		return s.handleMoveCommand(ctx, moveCmd)
	default:
		return nil, errors.New("unknown command")
	}
}

func (s *Service) handleMoveCommand(ctx context.Context, cmd parser.MoveCommand) (parser.Result, error) {
	room, err := s.FetchRoom(ctx, cmd.RoomId)
	if err != nil {
		return nil, err
	}

	result := parser.MoveResult{
		RoomId:    cmd.RoomId,
		Moved:     false,
		GameState: *room.GameState,
	}

	err = room.GameState.Board.MovePiece(cmd.FromPos, cmd.ToPos)
	if err == nil {
		result.Moved = true
		result.GameState = *room.GameState
	}

	return result, nil
}
