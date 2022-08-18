package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"

	"github.com/rs/zerolog"

	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/config"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/database/postgres"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/handlers"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/router"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/server"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/tasks"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/workers"
)

func main() {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	logger.Log().Msg("starting server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Log().Msg("starting parse configuration")

	cfg := config.New()

	logger.Log().Msg("ServerAddress: " + cfg.ServerAddress)
	logger.Log().Msg("AccrualSystemAddress: " + cfg.AccrualSystemAddress)
	logger.Log().Msg("DataBase: " + cfg.DataBase.DataBaseURI)

	db, err := sql.Open("postgres", cfg.DataBase.DataBaseURI)
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	logger.Log().Msg("finish db connection")

	repo := postgres.New(db, &logger)

	jobStore := postgres.NewJobStore(db)

	var listTask []tasks.TaskInterface
	listTask = append(listTask, tasks.NewCheckOrderStatusTask(cfg.AccrualSystemAddress, &logger, repo.ChangeOrderStatus))
	taskStore := tasks.NewTaskStore(listTask)

	wp := workers.New(jobStore, taskStore, &cfg.WorkerPool, &logger)

	go func() {
		wp.Run(ctx)
	}()

	logger.Log().Msg("starting setup db")

	_, err = postgres.RunMigration(cfg.DataBase.DataBaseURI)
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	logger.Log().Msg("finish setup db")

	h := handlers.New(repo, jobStore, &logger, cfg)

	router := router.New(h, cfg)

	s := server.New(ctx, router, cfg.ServerAddress)

	go func() error {
		err = s.Start()
		if err != nil {
			return err
		}

		logger.Log().Msgf("httpServer starting at: %s", cfg.ServerAddress)

		return nil
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}
}
