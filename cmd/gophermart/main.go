package main

import (
	"context"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/router"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.New(zerolog.DebugLevel)

	log.Log().Msg("Starting server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Log().Msg("Starting parse configuration")
	cfg := config.New()

	handler := router.New()

	s := server.New(ctx, handler, cfg)

	err := s.Start()
	if err != nil {
		panic(err)
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
