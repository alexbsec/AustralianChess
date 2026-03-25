package bot

import (
	"context"
	"log"
	"time"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/engine"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	"github.com/alexbsec/AustralianChess/backend/rooms"
)

const thinkDelay = 200

type TaipanBot struct {
	roomId      string
	botId       string
	color       chess.PieceColor
	difficulty  Difficulty
	roomService rooms.IService
	engine      engine.Engine
	stateCh     chan chess.GameState
	ctx         context.Context
	cancel      context.CancelFunc
	onMove      func(roomId string, result contracts.Result) error
}

func NewTaipan(roomId string, color chess.PieceColor, difficulty Difficulty, roomService rooms.IService, onMove func(roomId string, result contracts.Result) error) *TaipanBot {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaipanBot{
		roomId:      roomId,
		color:       color,
		difficulty:  difficulty,
		roomService: roomService,
		engine:      engine.NewChessEngine(),
		stateCh:     make(chan chess.GameState, 1),
		ctx:         ctx,
		cancel:      cancel,
		onMove:      onMove,
	}
}

func (b *TaipanBot) Start() {
	go b.loop()
}

func (b *TaipanBot) Stop() {
	b.cancel()
}

func (b *TaipanBot) NotifyState(state chess.GameState) {
	select {
	case b.stateCh <- state:
	default:
	}
}

func (b *TaipanBot) loop() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case state := <-b.stateCh:
			if !b.isMyTurn(state) {
				continue
			}
			if state.Result != nil {
				return
			}

			// Small delay so it doesn't feel robotic
			select {
			case <-time.After(thinkDelay * time.Microsecond):
			case <-b.ctx.Done():
				return
			}

			b.playMove(state)
		}
	}
}

func (b *TaipanBot) isMyTurn(state chess.GameState) bool {
	return state.Turn == b.color
}

func (b *TaipanBot) playMove(state chess.GameState) {
	move, ok := b.engine.IterativeDeepening(state.Board, b.color, int(b.difficulty))
	if !ok {
		log.Printf("taipan: no legal moves found in room %s", b.roomId)
		return
	}

	moveCmd := contracts.MoveCommand{
		RoomId:         b.roomId,
		PlayerId:       b.botId,
		RequesteeColor: b.color,
		FromPos:        move.FromPos,
		ToPos:          move.ToPos,
	}

	result, err := b.roomService.ExecuteCommand(b.ctx, moveCmd)
	if err != nil {
		log.Printf("taipan: failed to execute move in room %s: %v", b.roomId, err)
		return
	}

	if b.onMove != nil {
		b.onMove(b.roomId, result)
	}
}
