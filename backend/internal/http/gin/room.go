package gin

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/alexbsec/AustralianChess/backend/internal/chess"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/gin-gonic/gin"
)

func RoomHandler(
	ctx context.Context,
	routerGroup *gin.RouterGroup,
	roomService rooms.IService,
	wsHandler ws.IHandler,
) {
	routerGroup.Handle("GET", "/create", MakeNewRoom(roomService))
	routerGroup.Handle("GET", "/:id", ServeRoom(roomService))
	if wsHandler != nil {
		routerGroup.Handle("GET", "/ws/:id", Multiplayer(wsHandler))
		routerGroup.Handle("GET", "/ws/bot", Singleplayer(wsHandler))
	}
}

func MakeNewRoom(roomService rooms.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.MustGet("sessionId").(int64)
		_ = ctx.MustGet("userId").(int64)

		res, err := roomService.CreateRoom(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"roomId": res.RoomId})
	}
}

func ServeRoom(roomService rooms.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.MustGet("sessionId").(int64)
		_ = ctx.MustGet("userId").(int64)

		id := ctx.Param("id")
		res, err := roomService.DisplayRoom(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func Multiplayer(wsHandler ws.IHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomId := ctx.Param("id")
		if roomId == "" {
			log.Printf("missing room id in request")
			ctx.JSON(http.StatusBadRequest, ws.ErrorMessage{
				Type:    "error",
				Message: "missing room id",
			})
			return
		}

		wsHandler.HandleRoom(ctx, roomId)
	}
}

func Singleplayer(wsHandler ws.IHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomId := ctx.Query("roomId")
		if roomId == "" {
			ctx.JSON(http.StatusBadRequest, ws.ErrorMessage{
				Type:    "error",
				Message: "missing roomId query parameter",
			})
			return
		}

		difficultyInt, err := strconv.Atoi(ctx.Query("difficulty"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ws.ErrorMessage{
				Type:    "error",
				Message: "invalid difficulty query parameter",
			})
			return
		}

		playerPlayingAsInt, _ := strconv.Atoi(ctx.Query("player_playing_as"))

		playBotDTO := ws.PlayBotDTO{
			RoomId:          roomId,
			Difficulty:      bot.Difficulty(difficultyInt),
			PlayerPlayingAs: chess.PieceColor(playerPlayingAsInt),
		}

		wsHandler.PlayBot(ctx, playBotDTO)
	}
}
