package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
		token       string
	}

	tests := []struct {
		name        string
		query       string
		withoutBody bool
		body        string
		mockError   error
		mockUser    models.User
		want        want
	}{
		{
			name:        "the user has been successfully authenticated",
			query:       "/api/user/register",
			body:        `{"login": "login", "password": "12345"}`,
			withoutBody: false,
			mockError:   nil,
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid request format",
			query:       "/api/user/register",
			body:        ``,
			withoutBody: true,
			mockError:   errors.New("the body is missing"),
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid request body",
			query:       "/api/user/register",
			body:        `""`,
			withoutBody: true,
			mockError:   nil,
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "registration of a non-unique user",
			query:       "/api/user/register",
			body:        `{"login": "login", "password": "12345"}`,
			withoutBody: false,
			mockError:   NewErrorWithDB(errors.New("UniqConstraint"), "UniqConstraint"),
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusConflict,
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router.Post(tt.query, h.Register)

			if !tt.withoutBody {
				repoMock.EXPECT().CreateUser(gomock.Any(), tt.mockUser).Return(&tt.mockUser, tt.mockError)
			}

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
		})
	}
}

func TestHandlers_Login(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		token       string
	}

	tests := []struct {
		name        string
		query       string
		withoutBody bool
		body        string
		mockError   error
		mockUser    models.User
		want        want
	}{
		{
			name:        "successful user authorization",
			query:       "/api/user/login",
			body:        `{"login": "login", "password": "12345"}`,
			withoutBody: false,
			mockError:   nil,
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid request format",
			query:       "/api/user/register",
			body:        ``,
			withoutBody: true,
			mockError:   errors.New("the body is missing"),
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid request body",
			query:       "/api/user/register",
			body:        `""`,
			withoutBody: true,
			mockError:   nil,
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid user password",
			query:       "/api/user/register",
			body:        `{"login": "login", "password": "12345"}`,
			withoutBody: false,
			mockError:   errors.New("user not found"),
			mockUser: models.User{
				Login:    "login",
				Password: "12345",
			},
			want: want{
				code:        http.StatusConflict,
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router.Post(tt.query, h.Login)

			if !tt.withoutBody {
				repoMock.EXPECT().CheckPassword(gomock.Any(), tt.mockUser).Return(&tt.mockUser, tt.mockError)
			}

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestHandlers_CreateOrder(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		token       string
	}

	tests := []struct {
		name        string
		query       string
		withoutBody bool
		body        string
		mockError   error
		mockOrder   models.Order
		want        want
	}{
		{
			name:        "successful user authorization",
			query:       "/api/user/orders",
			body:        "79927398713",
			withoutBody: false,
			mockError:   nil,
			mockOrder: models.Order{
				UserID: "userID",
				Number: strconv.Itoa(79927398713),
				Status: "NEW",
			},
			want: want{
				code:        http.StatusAccepted,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router.Post(tt.query, h.CreateOrder)

			if !tt.withoutBody {
				repoMock.EXPECT().CreateOrder(gomock.Any(), tt.mockOrder).Return(nil)

				jobStoreMock.EXPECT().AddJob(gomock.Any(), gomock.Any()).Return(nil)
			}

			router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtx, "userID")))

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}
