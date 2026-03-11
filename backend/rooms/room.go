package rooms

import "context"

type IService interface {
	CreateRoom(ctx context.Context) (CreateRoomResponse, error)
	DisplayRoom(ctx context.Context, roomId string) (DisplayRoomResponse, error)
}
