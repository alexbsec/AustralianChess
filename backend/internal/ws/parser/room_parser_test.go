package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
)

func TestRoomParser_ParseMessage(t *testing.T) {
	parser := NewRoomParser()
	roomID := "room-123"

	tests := []struct {
		name        string
		input       string
		expectError bool
		validate    func(t *testing.T, cmd contracts.Command)
	}{
		{
			name: "valid make_move",
			input: `{
				"type": "make_move",
				"requestee_color": 0,
				"from_pos": { "row": 10, "col": 4 },
				"to_pos": { "row": 8, "col": 4 }
			}`,
			expectError: false,
			validate: func(t *testing.T, cmd contracts.Command) {
				moveCmd, ok := cmd.(contracts.MoveCommand)
				require.True(t, ok)

				require.Equal(t, roomID, moveCmd.RoomId)
				require.Equal(t, chess.PieceWhite, moveCmd.RequesteeColor)
				require.Equal(t, chess.Position{Row: 10, Col: 4}, moveCmd.FromPos)
				require.Equal(t, chess.Position{Row: 8, Col: 4}, moveCmd.ToPos)
			},
		},
		{
			name: "valid pawn_promote",
			input: `{
				"type": "pawn_promote",
				"requestee_color": 1,
				"promote_to": 7,
				"pawn_position": { "row": 1, "col": 10 },
				"destination_position": { "row": 0, "col": 10 }
			}`,
			expectError: false,
			validate: func(t *testing.T, cmd contracts.Command) {
				promoteCmd, ok := cmd.(contracts.PromoteCommand)
				require.True(t, ok)

				require.Equal(t, roomID, promoteCmd.RoomId)
				require.Equal(t, chess.PieceBlack, promoteCmd.RequesteeColor)
				require.Equal(t, chess.QueenPiece, promoteCmd.PromoteTo)

				require.Equal(t, chess.Position{Row: 1, Col: 10}, promoteCmd.PiecePosition)
				require.Equal(t, chess.Position{Row: 0, Col: 10}, promoteCmd.DestPosition)
			},
		},
		{
			name: "unknown command type",
			input: `{
				"type": "invalid_type"
			}`,
			expectError: true,
		},
		{
			name:        "invalid json",
			input:       `{invalid_json}`,
			expectError: true,
		},
		{
			name: "valid type but invalid inner json for make_move",
			input: `{
				"type": "make_move",
				"requestee_color": 0,
				"from_pos": "invalid",
				"to_pos": { "row": 8, "col": 4 }
			}`,
			expectError: true,
		},
		{
			name: "valid type but invalid inner json for promote",
			input: `{
				"type": "pawn_promote",
				"requestee_color": 1,
				"promote_to": "invalid",
				"pawn_position": { "row": 1, "col": 10 },
				"destination_position": { "row": 0, "col": 10 }
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.ParseMessage(roomID, []byte(tt.input))

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, cmd)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cmd)

			if tt.validate != nil {
				tt.validate(t, cmd)
			}
		})
	}
}
