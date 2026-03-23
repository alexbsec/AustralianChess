package parser

import (
	"encoding/json"
	"errors"

	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
)

type RoomParser struct {
}

func NewRoomParser() *RoomParser {
	return &RoomParser{}
}

func (d *RoomParser) ParseMessage(roomId string, rawMsg []byte) (contracts.Command, error) {
	var envelope Envelope
	if err := json.Unmarshal(rawMsg, &envelope); err != nil {
		return nil, err
	}

	switch envelope.Type {
	case "make_move":
		return d.parseMakeMoveCommand(roomId, rawMsg)
	case "pawn_promote":
		return d.parsePromoteCommand(roomId, rawMsg)
	default:
		return nil, errors.New("unknown command")
	}
}

func (d *RoomParser) parsePromoteCommand(roomId string, rawMsg []byte) (contracts.Command, error) {
	var msg PromotePawnMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return nil, err
	}

	promoteCmd := contracts.PromoteCommand{
		RoomId:         roomId,
		RequesteeColor: msg.RequesteeColor,
		PromoteTo:      msg.PromoteTo,
		PiecePosition:  msg.PawnPosition,
		DestPosition:   msg.DestPosition,
	}

	return promoteCmd, nil
}

func (d *RoomParser) parseMakeMoveCommand(roomId string, rawMsg []byte) (contracts.Command, error) {
	var msg MovePieceMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return nil, err
	}

	moveCmd := contracts.MoveCommand{
		RoomId:         roomId,
		RequesteeColor: msg.RequesteeColor,
		FromPos:        msg.FromPos,
		ToPos:          msg.ToPos,
	}

	return moveCmd, nil
}
