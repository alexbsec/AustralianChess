package bot

import (
	"context"
	"errors"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	engineMock "github.com/alexbsec/AustralianChess/backend/internal/chess/engine/mocks"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	roomMock "github.com/alexbsec/AustralianChess/backend/rooms/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestTaipanBot_playMove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMock.NewMockEngine(ctrl)
	mockRoom := roomMock.NewMockIService(ctrl) 

	ctx := context.Background()

	bot := &TaipanBot{
		engine:      mockEngine,
		roomService: mockRoom,
		ctx:         ctx,
		roomId:      "room-1",
		botId:       "bot-1",
		color:       chess.PieceWhite,
		difficulty:  2,
	}

	state := chess.GameState{
		Board: chess.Board{},
	}

	t.Run("no legal moves", func(t *testing.T) {
		mockEngine.EXPECT().
			IterativeDeepening(state.Board, bot.color, int(bot.difficulty)).
			Return(chess.Move{}, false)


		bot.playMove(state)
	})

	t.Run("execute command fails", func(t *testing.T) {
		move := chess.Move{
			FromPos: chess.Position{Row: 1, Col: 1},
			ToPos:   chess.Position{Row: 2, Col: 2},
		}

		mockEngine.EXPECT().
			IterativeDeepening(state.Board, bot.color, int(bot.difficulty)).
			Return(move, true)

		mockRoom.EXPECT().
			ExecuteCommand(ctx, gomock.Any()).
			Return(contracts.FailedCommandResult{}, errors.New("boom"))

		bot.playMove(state)
	})

	t.Run("success triggers onMove", func(t *testing.T) {
		move := chess.Move{
			FromPos: chess.Position{Col: 1, Row: 1},
			ToPos:   chess.Position{Col: 2, Row: 2},
		}

		expectedResult := contracts.MoveResult{
		}

		mockEngine.EXPECT().
			IterativeDeepening(state.Board, bot.color, int(bot.difficulty)).
			Return(move, true)

		mockRoom.EXPECT().
			ExecuteCommand(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, cmd contracts.MoveCommand) (contracts.Result, error) {
				// Assert command correctness inline
				require.Equal(t, "room-1", cmd.RoomId)
				require.Equal(t, "bot-1", cmd.PlayerId)
				require.Equal(t, chess.PieceWhite, cmd.RequesteeColor)
				require.Equal(t, move.FromPos, cmd.FromPos)
				require.Equal(t, move.ToPos, cmd.ToPos)

				return expectedResult, nil
			})

		called := false
		bot.onMove = func(roomId string, result contracts.Result) error {
			called = true
			require.Equal(t, "room-1", roomId)
			require.Equal(t, expectedResult, result)
			return nil
		}

		bot.playMove(state)

		require.True(t, called)
	})
}
