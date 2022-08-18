package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/config"
	"github.com/mkokoulin/go-musthave-diploma-tpl/internal/models"
)

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
		name      string
		query     string
		body      string
		mockError error
		mockOrder models.Order
		want      want
	}{
		{
			name:      "successful user creation",
			query:     "/api/user/orders",
			body:      "79927398713",
			mockError: nil,
			mockOrder: models.Order{
				UserID: "userID",
				Number: strconv.Itoa(79927398713),
				Status: "NEW",
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name:      "invalid order number",
			query:     "/api/user/orders",
			body:      "123456789",
			mockError: nil,
			mockOrder: models.Order{
				UserID: "userID",
				Number: strconv.Itoa(79927398713),
				Status: "NEW",
			},
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name:      "an order already registered by you",
			query:     "/api/user/orders",
			body:      "79927398713",
			mockError: NewErrorWithDB(errors.New("OrderAlreadyRegisterByYou"), "OrderAlreadyRegisterByYou"),
			mockOrder: models.Order{
				UserID: "userID",
				Number: strconv.Itoa(79927398713),
				Status: "NEW",
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:      "an order already registered",
			query:     "/api/user/orders",
			body:      "79927398713",
			mockError: NewErrorWithDB(errors.New("OrderAlreadyRegister"), "OrderAlreadyRegister"),
			mockOrder: models.Order{
				UserID: "userID",
				Number: strconv.Itoa(79927398713),
				Status: "NEW",
			},
			want: want{
				code: http.StatusConflict,
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

			repoMock.EXPECT().CreateOrder(gomock.Any(), tt.mockOrder).Return(tt.mockError).AnyTimes()

			jobStoreMock.EXPECT().AddJob(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtx, "userID")))

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestHandlers_GetOrders(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		token       string
	}

	tests := []struct {
		name       string
		query      string
		body       string
		mockError  error
		mockOrders []models.ResponseOrderWithAccrual
		want       want
	}{
		{
			name:      "successful receipt of orders",
			query:     "/api/user/orders",
			body:      "",
			mockError: nil,
			mockOrders: []models.ResponseOrderWithAccrual{
				{
					Accrual: 0,
					ResponseOrder: models.ResponseOrder{
						Number: "0",
						Status: "New",
					},
				},
			},
			want: want{
				code:        http.StatusOK,
				response:    `[{"number":"0","status":"New","uploaded_at":"0001-01-01T00:00:00Z"}]`,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:      "body no content",
			query:     "/api/user/orders",
			body:      "",
			mockError: nil,
			want: want{
				code:        http.StatusNoContent,
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router.Get(tt.query, h.GetOrders)

			repoMock.EXPECT().GetOrders(gomock.Any(), "userID").Return(tt.mockOrders, tt.mockError).AnyTimes()

			router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtx, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(body), "invalid response body")
		})
	}
}

func TestHandlers_GetBalance(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		token       string
	}

	tests := []struct {
		name        string
		query       string
		body        string
		mockError   error
		mockBalance models.UserBalance
		want        want
	}{
		{
			name:        "successful receipt of balance",
			query:       "/api/user/balance",
			mockBalance: models.UserBalance{},
			mockError:   nil,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json; charset=utf-8",
				response:    `{"current":0,"withdrawn":0}`,
			},
		},
		{
			name:  "unsuccessful receipt of balance",
			query: "/api/user/balance",
			mockBalance: models.UserBalance{
				Balance: 0,
				Spent:   0,
			},
			mockError: errors.New(""),
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
				response:    "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router.Get(tt.query, h.GetBalance)

			repoMock.EXPECT().GetBalance(gomock.Any(), "userID").Return(tt.mockBalance, tt.mockError).AnyTimes()

			router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtx, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(body), "invalid response body")
		})
	}
}

func TestHandlers_CreateWithdraw(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		token       string
	}

	tests := []struct {
		name         string
		query        string
		body         string
		mockError    error
		mockWithdraw models.Withdraw
		want         want
	}{
		{
			name:      "successful creating withdraw",
			query:     "/api/user/balance/withdraw",
			body:      `{"order":"79927398713","sum":0}`,
			mockError: nil,
			mockWithdraw: models.Withdraw{
				WithdrawOrder: models.WithdrawOrder{
					Order: "79927398713",
					Sum:   0,
				},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:      "successful creating withdraw",
			query:     "/api/user/balance/withdraw",
			body:      `{"order":"79927398713","sum":0}`,
			mockError: NewErrorWithDB(errors.New("NotEnoughBalanceForWithdraw"), "NotEnoughBalanceForWithdraw"),
			mockWithdraw: models.Withdraw{
				WithdrawOrder: models.WithdrawOrder{
					Order: "79927398713",
					Sum:   0,
				},
			},
			want: want{
				code: http.StatusPaymentRequired,
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

			router.Post(tt.query, h.CreateWithdraw)

			repoMock.EXPECT().CreateWithdraw(gomock.Any(), tt.mockWithdraw, "userID").Return(tt.mockError).AnyTimes()

			router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtx, "userID")))

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestHandlers_GetWithdrawals(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		token       string
	}

	tests := []struct {
		name          string
		query         string
		body          string
		mockError     error
		mockWithdraws []models.WithdrawOrder
		want          want
	}{
		{
			name:          "successful receipt withdrawals",
			query:         "/api/user/balance/withdrawals",
			mockError:     nil,
			mockWithdraws: []models.WithdrawOrder{models.WithdrawOrder{}},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				response:    `[{"order":"","sum":0,"processed_at":"0001-01-01T00:00:00Z"}]`,
			},
		},
		{
			name:          "no content",
			query:         "/api/user/balance/withdrawals",
			mockError:     nil,
			mockWithdraws: []models.WithdrawOrder{},
			want: want{
				code:     http.StatusNoContent,
				response: ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, tt.query, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
			cfg := config.New()

			repoMock := NewMockRepository(ctrl)
			jobStoreMock := NewMockJobStore(ctrl)

			h := New(repoMock, jobStoreMock, &logger, cfg)

			router.Get(tt.query, h.GetWithdrawals)

			repoMock.EXPECT().GetWithdrawals(gomock.Any(), "userID").Return(tt.mockWithdraws, tt.mockError).AnyTimes()

			router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtx, "userID")))

			response := w.Result()

			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(body), "invalid response body")
		})
	}
}
