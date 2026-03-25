package bot

import (
	"fmt"
	"sync"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	"github.com/alexbsec/AustralianChess/backend/rooms"
)

type BotManager struct {
	mtx  sync.RWMutex
	bots map[string]Bot
}

func NewBotManager() *BotManager {
	return &BotManager{
		bots: make(map[string]Bot),
	}
}

func (m *BotManager) SpawnBot(roomId string, color chess.PieceColor, difficulty Difficulty, roomService rooms.IService, onMove func(string, contracts.Result) error) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, exists := m.bots[roomId]; exists {
		return fmt.Errorf("bot already exists in room %s", roomId)
	}

	taipan := NewTaipan(roomId, color, difficulty, roomService, onMove)
	m.bots[roomId] = taipan
	taipan.Start()
	return nil
}

func (m *BotManager) NotifyState(roomId string, state chess.GameState) {
	m.mtx.RLock()
	taipan, ok := m.bots[roomId]
	m.mtx.RUnlock()
	if ok {
		taipan.NotifyState(state)	
	}
}

func (m *BotManager) RemoveBot(roomId string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if taipan, ok := m.bots[roomId]; ok {
		taipan.Stop()
		delete(m.bots, roomId)
	}
}

func (m *BotManager) HasBot(roomId string) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	_, ok := m.bots[roomId]
	return ok
}
