package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/rs/zerolog"
)

func main() {
	logger := logger.New(zerolog.DebugLevel)

	logger.Log("Starting server")

	//ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger.Log("Starting parse configuration")

	cfg := config.New()

	logger.Log(cfg.Token.AccessTokenSecret)

	_, err := sql.Open("postgres", cfg.DataBase.DataBaseURI)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Log("Finish db connection")

	//repo := postgres.New(db)
	//
	//logger.Log("Starting setup db")
	//
	////_, err = database.RunMigration(cfg.DataBase.DataBaseURI)
	////if err != nil {
	////	logger.Fatal(err.Error())
	////}
	//
	//logger.Log("Finish setup db")
	//
	//handlers := handlers.New(repo, cfg)
	//
	//router := router.New(handlers, cfg)
	//
	//log.Println(cfg)
	//
	//s := server.New(ctx, router, cfg.ServerAddress)
	//
	//g, ctx := errgroup.WithContext(ctx)
	//
	//go func() error {
	//	err = s.Start()
	//	if err != nil {
	//		return err
	//	}
	//
	//	logger.Log("httpServer starting at: " + cfg.ServerAddress)
	//
	//	return nil
	//}()
	//
	//select {
	//case <-interrupt:
	//	cancel()
	//	break
	//case <-ctx.Done():
	//	break
	//}
	//
	//err = g.Wait()
	//if err != nil {
	//	logger.Log("server returning an error: " + err.Error())
	//	os.Exit(2)
	//}
}
