package game

import (
	"context"
	"fmt"
	"log"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	wsTypes "github.com/alexbsec/AustralianChess/backend/internal/ws/ws_types"
	"github.com/alexbsec/AustralianChess/backend/rooms"
)

type GameOrchestrator struct {
	roomService    rooms.IService
	botManager     bot.Manager
	onMoveCallback func(string, contracts.Result) error
}

func NewGameOrchestrator(roomService rooms.IService, botManager bot.Manager, callback func(string, contracts.Result) error) *GameOrchestrator {
	return &GameOrchestrator{
		roomService:    roomService,
		botManager:     botManager,
		onMoveCallback: callback,
	}
}

func (g *GameOrchestrator) HandlePlayBotRoom(ctx context.Context, roomId string, playerColor chess.PieceColor, difficulty bot.Difficulty) error {
	room, err := g.roomService.FetchRoom(ctx, roomId)
	if err != nil {
		return err
	}

	botColor := chess.OponentColor(playerColor)
	botId := bot.DifficultyToBotId(difficulty)
	err = g.botManager.SpawnBot(room.Id, botColor, difficulty, g.roomService, g.onMoveCallback)
	if err != nil {
		_ = g.roomService.DeleteRoom(ctx, room.Id)
		return fmt.Errorf("failed to spawn bot: %w", err)
	}

	_, err = g.roomService.UpdatePlayerJoined(ctx, roomId, botId, botColor)
	if err != nil {
		_ = g.roomService.DeleteRoom(ctx, room.Id)
		return fmt.Errorf("failed to update bot joined: %w", err)
	}

	return nil
}

func (g *GameOrchestrator) HandleConnect(ctx context.Context, roomId string, client *wsTypes.Client) (contracts.Result, chess.GameState, error) {
	room, err := g.roomService.FetchRoom(ctx, roomId)
	if err != nil {
		log.Printf("room %s not found: %v", roomId, err)
		return nil, chess.GameState{}, err
	}

	if client == nil {
		return nil, chess.GameState{}, ErrNilClient
	}

	// Player joins game
	if client.Role == wsTypes.PlayerOne || client.Role == wsTypes.PlayerTwo {
		result, err := g.roomService.UpdatePlayerJoined(
			ctx,
			roomId,
			client.PlayerId,
			*client.Color,
		)
		if err != nil {
			return nil, chess.GameState{}, err
		}

		// Notify bot of initial state so it can move first if it's white
		if g.botManager.HasBot(roomId) {
			g.notifyBotAfterJoin(ctx, roomId, result)
		}

		return result, *room.GameState, nil
	}

	// spectator
	return nil, *room.GameState, nil
}

func (g *GameOrchestrator) HandleCommand(ctx context.Context, roomId string, client *wsTypes.Client, cmd contracts.Command) (contracts.Result, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	cmd = cmd.SetPlayerId(client.PlayerId)
	cmd = cmd.SetColor(*client.Color)

	result, err := g.roomService.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	if g.botManager.HasBot(roomId) {
		if state, ok := gameStateFromResult(result); ok {
			g.botManager.NotifyState(roomId, state)
		}
	}
	return result, nil
}

func (g *GameOrchestrator) HandleDisconnect(roomId string, client *wsTypes.Client) {
	if client == nil {
		return
	}

	log.Printf("orchestrator: player %s disconnected from room %s", client.PlayerId, roomId)

	// if a bot is playing
	if g.botManager.HasBot(roomId) {
		g.botManager.RemoveBot(roomId)
	}
}

func (g *GameOrchestrator) HandleDestroyRoom(ctx context.Context, roomId string) error {
	return g.roomService.DeleteRoom(ctx, roomId)
}

func (g *GameOrchestrator) notifyBotAfterJoin(ctx context.Context, roomId string, result contracts.Result) {
	joinResult, ok := result.(contracts.PlayerJoinedResult)
	if !ok || !joinResult.GameStarted {
		return
	}
	room, err := g.roomService.FetchRoom(ctx, roomId)
	if err != nil || room.GameState == nil {
		return
	}
	g.botManager.NotifyState(roomId, *room.GameState)
}

func gameStateFromResult(result contracts.Result) (chess.GameState, bool) {
	switch r := result.(type) {
	case contracts.MoveResult:
		return r.GameState, r.Moved
	case contracts.PromotionResult:
		return r.GameState, r.Promoted
	default:
		return chess.GameState{}, false
	}
}
