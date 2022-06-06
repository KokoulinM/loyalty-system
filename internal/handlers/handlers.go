package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/auth"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
}

type Handlers struct {
	repo Repository
	cfg  config.Config
}

type ErrorWithDB struct {
	Err   error
	Title string
}

func (err *ErrorWithDB) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *ErrorWithDB) Unwrap() error {
	return err.Err
}

func NewErrorWithDB(err error, title string) error {
	return &ErrorWithDB{
		Err:   err,
		Title: title,
	}
}

func New(repo Repository, cfg config.Config) *Handlers {
	return &Handlers{
		repo: repo,
		cfg:  cfg,
	}
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Add("Content-Type", "application/json; charset=utf-8")

	user := models.User{}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if body != nil {
		http.Error(w, "the body is missing", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUser, err := h.repo.CreateUser(r.Context(), user)
	var dbErr *ErrorWithDB

	if errors.As(err, &dbErr) && dbErr.Title == "UniqConstraint" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	token, err := auth.CreateToken(newUser.ID, h.cfg.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", "Bearer "+token.AccessToken)

	w.WriteHeader(http.StatusOK)
}
