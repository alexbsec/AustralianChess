package parser

type IParser interface {
	ParseMessage(roomId string, rawMsg []byte) (Command, error)
}
