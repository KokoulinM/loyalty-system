package router

import (
	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/handlers/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(h *handlers.Handlers, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		router.Post("/api/user/register", h.Register)
		router.Post("/api/user/login", h.Login)
		router.Route("/api/user/orders", func(r chi.Router) {
			r.Use(middlewares.JWTMiddleware(&cfg.Token))
			r.Post("/", h.CreateOrder)
			r.Get("/", h.GetOrders)
		})
	})

	return router
}
