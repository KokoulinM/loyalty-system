package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/app/handlers"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/auth"
)

func JWTMiddleware(cfg *config.ConfigToken) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "register") && !strings.Contains(r.URL.Path, "login") {
				token, err := auth.ValidateToken(r, cfg)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				userID := token.Claims.(jwt.MapClaims)["user_id"]

				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), handlers.UserIDCtx, userID)))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
