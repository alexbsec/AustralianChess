package rooms_test

import (
	"context"
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	contractMocks "github.com/alexbsec/AustralianChess/backend/internal/contracts/mocks"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/golang/mock/gomock"
)

func TestExecuteCommand_NilCommand(t *testing.T) {
	svc := rooms.NewService()

	_, err := svc.ExecuteCommand(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil command")
	}
}

func TestExecuteCommand_RoomNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := rooms.NewService()

	cmd := contractMocks.NewMockCommand(ctrl)
	cmd.EXPECT().GetRoomId().Return("invalid")

	_, err := svc.ExecuteCommand(context.Background(), cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExecuteCommand_GameNotStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := rooms.NewService()
	roomRes, _ := svc.CreateRoom(context.Background())

	cmd := contractMocks.NewMockCommand(ctrl)
	cmd.EXPECT().GetRoomId().Return(roomRes.RoomId)

	result, err := svc.ExecuteCommand(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	failed, ok := result.(contracts.FailedCommandResult)
	if !ok {
		t.Fatal("expected FailedCommandResult")
	}

	if failed.GameStarted {
		t.Fatal("game should not be started")
	}
}

func TestExecuteCommand_UnknownCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := rooms.NewService()
	roomRes, _ := svc.CreateRoom(context.Background())

	// start game
	_, _ = svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p1", chess.PieceWhite)
	_, _ = svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p2", chess.PieceBlack)

	cmd := contractMocks.NewMockCommand(ctrl)
	cmd.EXPECT().GetRoomId().Return(roomRes.RoomId)

	_, err := svc.ExecuteCommand(context.Background(), cmd)
	if err == nil {
		t.Fatal("expected unknown command error")
	}
}

func TestUpdatePlayerJoined_StartsGame(t *testing.T) {
	svc := rooms.NewService()
	roomRes, _ := svc.CreateRoom(context.Background())

	_, err := svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p1", chess.PieceWhite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p2", chess.PieceBlack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	join, ok := result.(contracts.PlayerJoinedResult)
	if !ok {
		t.Fatalf("expected PlayerJoinedResult, got %T", result)
	}

	if !join.GameStarted {
		t.Fatal("expected game to be started")
	}

	if !join.Success {
		t.Fatal("expected success")
	}
}

func TestCreateRoom_AndFetchRoom(t *testing.T) {
	svc := rooms.NewService()

	res, err := svc.CreateRoom(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	room, err := svc.FetchRoom(context.Background(), res.RoomId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if room.Id != res.RoomId {
		t.Fatal("room mismatch")
	}
}

func TestDeleteRoom(t *testing.T) {
	svc := rooms.NewService()
	res, _ := svc.CreateRoom(context.Background())

	err := svc.DeleteRoom(context.Background(), res.RoomId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.FetchRoom(context.Background(), res.RoomId)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestDisplayRoom_InvalidRoom(t *testing.T) {
	svc := rooms.NewService()

	_, err := svc.DisplayRoom(context.Background(), "invalid")
	if err == nil {
		t.Fatal("expected error for invalid room")
	}
}


func TestFetchRoom_Invalid(t *testing.T) {
	svc := rooms.NewService()

	_, err := svc.FetchRoom(context.Background(), "invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}


func TestDeleteRoom_Invalid(t *testing.T) {
	svc := rooms.NewService()

	err := svc.DeleteRoom(context.Background(), "invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}


func TestUpdatePlayerJoined_InvalidRoom(t *testing.T) {
	svc := rooms.NewService()

	res, err := svc.UpdatePlayerJoined(context.Background(), "invalid", "p1", chess.PieceWhite)
	if err == nil {
		t.Fatal("expected error")
	}

	join := res.(contracts.PlayerJoinedResult)
	if join.Success {
		t.Fatal("should not be successful")
	}
}


func TestUpdatePlayerJoined_OnePlayerOnly(t *testing.T) {
	svc := rooms.NewService()
	roomRes, _ := svc.CreateRoom(context.Background())

	res, err := svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p1", chess.PieceWhite)
	if err != nil {
		t.Fatal(err)
	}

	join := res.(contracts.PlayerJoinedResult)

	if join.GameStarted {
		t.Fatal("game should not start with one player")
	}
}


func TestUpdatePlayerJoined_OverwriteSameColor(t *testing.T) {
	svc := rooms.NewService()
	roomRes, _ := svc.CreateRoom(context.Background())

	_, _ = svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p1", chess.PieceWhite)
	res, _ := svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p2", chess.PieceWhite)

	join := res.(contracts.PlayerJoinedResult)

	if !join.Success {
		t.Fatal("expected success even if overwriting")
	}

	if join.GameStarted {
		t.Fatal("game should not start with only one color assigned")
	}
}


func TestExecuteCommand_StartedGame_ReturnsUnknownCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := rooms.NewService()
	roomRes, _ := svc.CreateRoom(context.Background())

	_, _ = svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p1", chess.PieceWhite)
	_, _ = svc.UpdatePlayerJoined(context.Background(), roomRes.RoomId, "p2", chess.PieceBlack)

	cmd := contractMocks.NewMockCommand(ctrl)
	cmd.EXPECT().GetRoomId().Return(roomRes.RoomId)

	res, err := svc.ExecuteCommand(context.Background(), cmd)

	if err == nil {
		t.Fatal("expected error")
	}

	if res != nil {
		t.Fatal("expected nil result on error")
	}
}
