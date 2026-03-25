package bot_test

import (
	"testing"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	roomMocks "github.com/alexbsec/AustralianChess/backend/rooms/mocks"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewBotManager(t *testing.T) {
	m := bot.NewBotManager()
	require.NotNil(t, m)
}

func TestBotManager_HasBot_False(t *testing.T) {
	m := bot.NewBotManager()
	require.False(t, m.HasBot("room-1"))
}

func TestBotManager_SpawnBot_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := roomMocks.NewMockIService(ctrl)
	m := bot.NewBotManager()

	err := m.SpawnBot("room-1", chess.PieceBlack, bot.EasyDifficulty, mockService, nil)
	require.NoError(t, err)
	require.True(t, m.HasBot("room-1"))

	// Cleanup
	m.RemoveBot("room-1")
}

func TestBotManager_SpawnBot_Duplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := roomMocks.NewMockIService(ctrl)
	m := bot.NewBotManager()

	err := m.SpawnBot("room-1", chess.PieceBlack, bot.EasyDifficulty, mockService, nil)
	require.NoError(t, err)

	// Spawn again in same room
	err = m.SpawnBot("room-1", chess.PieceBlack, bot.EasyDifficulty, mockService, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bot already exists")

	m.RemoveBot("room-1")
}

func TestBotManager_RemoveBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := roomMocks.NewMockIService(ctrl)
	m := bot.NewBotManager()

	err := m.SpawnBot("room-2", chess.PieceBlack, bot.EasyDifficulty, mockService, nil)
	require.NoError(t, err)
	require.True(t, m.HasBot("room-2"))

	m.RemoveBot("room-2")
	require.False(t, m.HasBot("room-2"))
}

func TestBotManager_RemoveBot_NonExistent(t *testing.T) {
	m := bot.NewBotManager()
	// Should not panic
	m.RemoveBot("nonexistent")
}

func TestBotManager_NotifyState_NoBot(t *testing.T) {
	m := bot.NewBotManager()
	// Should not panic when bot doesn't exist
	state := chess.GameState{
		Board: chess.NewInitialBoard(),
		Turn:  chess.PieceWhite,
	}
	m.NotifyState("nonexistent", state)
}

func TestBotManager_NotifyState_WithBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := roomMocks.NewMockIService(ctrl)
	m := bot.NewBotManager()

	err := m.SpawnBot("room-3", chess.PieceBlack, bot.EasyDifficulty, mockService, nil)
	require.NoError(t, err)

	// Notify with wrong turn so bot won't actually try to move
	state := chess.GameState{
		Board: chess.NewInitialBoard(),
		Turn:  chess.PieceWhite, // bot is black, so it won't move
	}
	m.NotifyState("room-3", state)

	// Give a moment for goroutine
	time.Sleep(10 * time.Millisecond)

	m.RemoveBot("room-3")
}

func TestBotManager_MultipleBots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := roomMocks.NewMockIService(ctrl)
	m := bot.NewBotManager()

	onMove := func(roomId string, result contracts.Result) error { return nil }

	err1 := m.SpawnBot("room-a", chess.PieceBlack, bot.EasyDifficulty, mockService, onMove)
	err2 := m.SpawnBot("room-b", chess.PieceBlack, bot.MediumDifficulty, mockService, onMove)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.True(t, m.HasBot("room-a"))
	require.True(t, m.HasBot("room-b"))

	m.RemoveBot("room-a")
	m.RemoveBot("room-b")
}
