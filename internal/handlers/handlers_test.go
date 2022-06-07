package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/router"
)

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
			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

			cfg := config.New()
			ctx, _ := context.WithCancel(context.Background())
			repoMock := new(MockRepository)

			h := New(repoMock, jobStore, &logger, cfg)

			router := router.New(h, cfg)

			repoMock.CreateUser(ctx, tt.mockUser)

			w := httptest.NewRecorder()
			body := strings.NewReader(tt.body)
			req, _ := http.NewRequest(http.MethodPost, tt.query, body)
			router.ServeHTTP(w, req)
		})
	}
}
