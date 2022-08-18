package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/config"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/handlers"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/middlewares"
)

func New(h *handlers.Handlers, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/", func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(&cfg.Token), middlewares.GzipMiddleware)
		r.Post("/api/user/register", h.Register)
		r.Post("/api/user/login", h.Login)
		r.Post("/api/user/orders", h.CreateOrder)
		r.Get("/api/user/orders", h.GetOrders)
		r.Get("/api/user/balance", h.GetBalance)
		r.Post("/api/user/balance/withdraw", h.CreateWithdraw)
		r.Get("/api/user/balance/withdrawals", h.GetWithdrawals)
	})

	return router
}
