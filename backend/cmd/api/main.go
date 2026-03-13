package main

import (
	"context"
	"log"

	"github.com/alexbsec/AustralianChess/backend/internal/config"
	"github.com/alexbsec/AustralianChess/backend/internal/db"
	"github.com/alexbsec/AustralianChess/backend/internal/db/repositories"
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

	roomSvc := rooms.NewService()
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(roomSvc, hub)

	router := ginChess.MakeHandlers(ctx, roomSvc, userSvc, wsHandler)

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
