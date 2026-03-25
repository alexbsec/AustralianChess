package main

import (
	"context"
	"log"

	"github.com/alexbsec/AustralianChess/backend/internal/auth"
	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/alexbsec/AustralianChess/backend/internal/config"
	"github.com/alexbsec/AustralianChess/backend/internal/contracts"
	"github.com/alexbsec/AustralianChess/backend/internal/db"
	"github.com/alexbsec/AustralianChess/backend/internal/db/repositories"
	"github.com/alexbsec/AustralianChess/backend/internal/game"
	ginChess "github.com/alexbsec/AustralianChess/backend/internal/http/gin"
	"github.com/alexbsec/AustralianChess/backend/internal/ws"
	"github.com/alexbsec/AustralianChess/backend/rooms"
	"github.com/alexbsec/AustralianChess/backend/sessions"
	"github.com/alexbsec/AustralianChess/backend/users"
)

func main() {
	ctx := context.Background()

	envs := config.LoadEnvs()
	database, err := connectDatabase(ctx, envs)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)

	sessionSvc := sessions.NewService(sessionRepo)
	userSvc := users.NewService(userRepo, sessionSvc)

	hub := ws.NewHub()
	roomSvc := rooms.NewService()
	botMangaer := bot.NewBotManager()
	gameOrchestrator := game.NewGameOrchestrator(roomSvc, botMangaer, func(roomId string, result contracts.Result) error {
		return hub.Broadcast(roomId, result, true)
	})

	wsHandler := ws.NewHandler(userSvc, hub, gameOrchestrator)
	authorizer := auth.NewAuthorizer(sessionSvc)

	router := ginChess.MakeHandlers(ctx, roomSvc, userSvc, wsHandler, authorizer)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func connectDatabase(ctx context.Context, envs *config.Environments) (db.Database, error) {
	connString := db.PostgresConnStringFromEnv(map[string]any{
		"DB_HOST":     envs.DBHost,
		"DB_PORT":     envs.DBPort,
		"DB_USERNAME": envs.DBUsername,
		"DB_PASSWORD": envs.DBPassword,
		"DB_NAME":     envs.DBName,
	})

	database := db.NewPostgresDB(connString)
	if err := database.Connect(ctx); err != nil {
		return nil, err
	}

	return database, nil
}
