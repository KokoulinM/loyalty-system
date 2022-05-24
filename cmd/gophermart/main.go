package main

import (
	"context"
	"database/sql"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/migratons"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/database/postgres"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/router"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/server"
	"github.com/rs/zerolog"
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
	migratons.Migrations(db, logger)

	logger.Log("Finish setup db")

	handlers := handlers.New(repo, cfg)

	router := router.New(handlers, cfg)

	s := server.New(ctx, router, cfg)

	err = s.Start()
	if err != nil {
		logger.Fatal(err.Error())
	}

	//token, err := auth.CreateToken("1", "", &cfg.Token)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for {
	//	jwt, err := auth.ValidateToken(token.AccessToken, &cfg.Token)
	//	if err != nil {
	//		token, err = auth.RefreshToken(token.RefreshToken, &cfg.Token)
	//		if err != nil {
	//			panic(err)
	//		}
	//	} else {
	//		fmt.Println(jwt.Raw)
	//	}
	//
	//	time.Sleep(1 * time.Second)
	//}
}
