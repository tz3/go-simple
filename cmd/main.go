package main

import (
	_ "github.com/lib/pq" // not used directly
	"go.uber.org/zap"

	"github.com/tz3/go-simple/config"
	"github.com/tz3/go-simple/database"
	"github.com/tz3/go-simple/server"
	"github.com/tz3/go-simple/user/handler"
	"github.com/tz3/go-simple/user/repository"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	db, err := database.GetDB(log)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(log, userRepo)

	srv := server.NewServer(config.DefaultPort)
	srv.RegisterHandler("/users", userHandler.UserHandler)
	srv.Start()
}
