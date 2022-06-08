package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
)

func newRouter(h *Handlers, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/", func(r chi.Router) {
		//r.Use(middlewares.JWTMiddleware(&cfg.Token), middlewares.GzipMiddleware)
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

func TestHandlers_Register(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name     string
		query    string
		body     string
		mockUser models.User
		want     want
	}{
		{
			name:  "пользователь успешно аутентифицирован",
			query: "/api/user/register",
			body:  `{"login": "login", "password": "12345"}`,
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        200,
				contentType: "application/json; charset=utf-8",
			},
		},
		//{
		//	name:  "неверный формат запроса",
		//	query: "/api/user/register",
		//},
		//{
		//	name:  "неверная пара логин/пароль",
		//	query: "/api/user/register",
		//},
		//{
		//	name:  "внутренняя ошибка сервера",
		//	query: "/api/user/register",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router := newRouter(h, cfg)

			//repoMock.CreateUser(ctx, tt.mockUser)

			w := httptest.NewRecorder()
			body := strings.NewReader(tt.body)
			req, _ := http.NewRequest(http.MethodPost, tt.query, body)
			router.ServeHTTP(w, req)

			_, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			//for _, field := range tt.want.InResponse {
			//	assert.Contains(t, string(resBody), field)
			//}
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header()["Content-Type"][0])
		})
	}
}
