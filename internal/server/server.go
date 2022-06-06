package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

type server struct {
	srv *http.Server
	ctx context.Context
}

func New(ctx context.Context, handler *chi.Mux, addr string) *server {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &server{
		srv: s,
		ctx: ctx,
	}
}

func (s *server) Start() error {
	err := s.srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *server) Stop() error {
	err := s.srv.Shutdown(s.ctx)
	if err != nil {
		return err
	}

	return nil
}
