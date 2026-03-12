package parser

import (
	"encoding/json"
	"errors"
)

type RoomParser struct {
}

func NewRoomParser() *RoomParser {
	return &RoomParser{}
}

func (d *RoomParser) ParseMessage(roomId string, rawMsg []byte) (Command, error) {
	var envelope Envelope
	if err := json.Unmarshal(rawMsg, &envelope); err != nil {
		return nil, err
	}

	switch envelope.Type {
	case "make_move":
		return d.parseMakeMoveCommand(roomId, rawMsg)
	default:
		return nil, errors.New("unknown command")
	}
}

func (d *RoomParser) parseMakeMoveCommand(roomId string, rawMsg []byte) (Command, error) {
	var msg MovePieceMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return nil, err
	}

	moveCmd := MoveCommand{
		RoomId:         roomId,
		RequesteeColor: msg.RequesteeColor,
		FromPos:        msg.FromPos,
		ToPos:          msg.ToPos,
	}

	return moveCmd, nil
}
