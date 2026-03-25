package contracts_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	"github.com/stretchr/testify/require"
)

func TestMoveCommand_GetRoomId(t *testing.T) {
	cmd := contracts.MoveCommand{RoomId: "room-1"}
	require.Equal(t, "room-1", cmd.GetRoomId())
}

func TestMoveCommand_SetPlayerId(t *testing.T) {
	cmd := contracts.MoveCommand{RoomId: "room-1"}
	updated := cmd.SetPlayerId("player-1")
	// SetPlayerId returns a new Command (value receiver)
	require.NotNil(t, updated)
}

func TestMoveCommand_SetColor(t *testing.T) {
	cmd := contracts.MoveCommand{RoomId: "room-1"}
	updated := cmd.SetColor(chess.PieceBlack)
	require.NotNil(t, updated)
}

func TestPromoteCommand_GetRoomId(t *testing.T) {
	cmd := contracts.PromoteCommand{RoomId: "room-2"}
	require.Equal(t, "room-2", cmd.GetRoomId())
}

func TestPromoteCommand_SetPlayerId(t *testing.T) {
	cmd := contracts.PromoteCommand{RoomId: "room-2"}
	updated := cmd.SetPlayerId("player-1")
	require.NotNil(t, updated)
}

func TestPromoteCommand_SetColor(t *testing.T) {
	cmd := contracts.PromoteCommand{RoomId: "room-2"}
	updated := cmd.SetColor(chess.PieceWhite)
	require.NotNil(t, updated)
}
