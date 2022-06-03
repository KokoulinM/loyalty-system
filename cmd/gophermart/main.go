package main

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/database"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/database/postgres"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/router"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/server"
)

func main() {
	logger := logger.New(zerolog.DebugLevel)

	logger.Log("Starting server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Log("Starting parse configuration")

	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DataBaseURI)
	logger.Log("Finish db connection")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer db.Close()

	repo := postgres.New(db)

	logger.Log("Starting setup db")
	database.Migrations(db, logger)

	logger.Log("Finish setup db")

	handlers := handlers.New(repo, &cfg)

	router := router.New(handlers, &cfg)

	s := server.New(ctx, router, cfg.ServerAddress)

	err = s.Start()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
