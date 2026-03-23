package parser

import "github.com/alexbsec/AustralianChess/backend/internal/contracts"

type IParser interface {
	ParseMessage(roomId string, rawMsg []byte) (contracts.Command, error)
}
