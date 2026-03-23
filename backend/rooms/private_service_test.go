package rooms

import (
	"context"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/engine"
	engineMocks "github.com/alexbsec/AustralianChess/backend/internal/chess/engine/mocks"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	"github.com/golang/mock/gomock"
)

func TestHandleMoveCommand_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.MoveCommand{
		RoomId:         roomRes.RoomId,
		RequesteeColor: chess.PieceWhite,
		FromPos:        chess.Position{Row: 1, Col: 1},
		ToPos:          chess.Position{Row: 2, Col: 1},
	}

	mockEngine.
		EXPECT().
		ValidateAndMove(room.GameState, cmd.RequesteeColor, cmd.FromPos, cmd.ToPos).
		Return(engine.ValidationResponse{
			Type:    engine.SuccessResponse,
			Message: "ok",
		})

	mockEngine.
		EXPECT().
		GameResult(room.GameState, chess.OponentColor(cmd.RequesteeColor)).
		Return(nil)

	res, err := svc.handleMoveCommand(context.Background(), room, cmd)
	if err != nil {
		t.Fatal(err)
	}

	move := res.(contracts.MoveResult)

	if !move.Moved {
		t.Fatal("expected move success")
	}
}

func TestHandleMoveCommand_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.MoveCommand{
		RoomId:         roomRes.RoomId,
		RequesteeColor: chess.PieceWhite,
	}

	mockEngine.
		EXPECT().
		ValidateAndMove(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(engine.ValidationResponse{
			Type:    engine.FailureResponse,
			Message: "invalid move",
		})

	res, err := svc.handleMoveCommand(context.Background(), room, cmd)
	if err != nil {
		t.Fatal(err)
	}

	move := res.(contracts.MoveResult)

	if move.Moved {
		t.Fatal("expected move failure")
	}
}

func TestHandleMoveCommand_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.MoveCommand{
		RoomId: roomRes.RoomId,
	}

	mockEngine.
		EXPECT().
		ValidateAndMove(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(engine.ValidationResponse{
			Type:    engine.ErrorResponse,
			Message: "boom",
		})

	_, err := svc.handleMoveCommand(context.Background(), room, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHandleMoveCommand_CheckmateEndReason(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.MoveCommand{
		RoomId:         roomRes.RoomId,
		RequesteeColor: chess.PieceWhite,
	}

	result := chess.ResultCheckmate

	mockEngine.
		EXPECT().
		ValidateAndMove(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(engine.ValidationResponse{
			Type: engine.SuccessResponse,
		})

	mockEngine.
		EXPECT().
		GameResult(gomock.Any(), gomock.Any()).
		Return(&result)

	res, err := svc.handleMoveCommand(context.Background(), room, cmd)
	if err != nil {
		t.Fatal(err)
	}

	move := res.(contracts.MoveResult)

	if move.GameState.EndReason == nil {
		t.Fatal("expected checkmate end reason")
	}
}

func TestHandlePromoteCommand_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.PromoteCommand{
		RoomId:         roomRes.RoomId,
		RequesteeColor: chess.PieceWhite,
	}

	mockEngine.
		EXPECT().
		PromotePawn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(engine.ValidationResponse{
			Type: engine.SuccessResponse,
		})

	res, err := svc.handlePromoteCommand(context.Background(), room, cmd)
	if err != nil {
		t.Fatal(err)
	}

	promote := res.(contracts.PromotionResult)

	if !promote.Promoted {
		t.Fatal("expected promotion success")
	}
}

func TestHandlePromoteCommand_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.PromoteCommand{
		RoomId: roomRes.RoomId,
	}

	mockEngine.
		EXPECT().
		PromotePawn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(engine.ValidationResponse{
			Type: engine.FailureResponse,
		})

	res, err := svc.handlePromoteCommand(context.Background(), room, cmd)
	if err != nil {
		t.Fatal(err)
	}

	promote := res.(contracts.PromotionResult)

	if promote.Promoted {
		t.Fatal("expected promotion failure")
	}
}

func TestHandlePromoteCommand_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := engineMocks.NewMockEngine(ctrl)
	svc := NewServiceWithEngine(mockEngine)

	roomRes, _ := svc.CreateRoom(context.Background())
	room, _ := svc.FetchRoom(context.Background(), roomRes.RoomId)

	cmd := contracts.PromoteCommand{
		RoomId: roomRes.RoomId,
	}

	mockEngine.
		EXPECT().
		PromotePawn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(engine.ValidationResponse{
			Type:    engine.ErrorResponse,
			Message: "fail",
		})

	_, err := svc.handlePromoteCommand(context.Background(), room, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}
