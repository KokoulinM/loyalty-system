package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/golang-jwt/jwt"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID string, email string, cfg *config.ConfigToken) (*TokenDetails, error) {
	td := &TokenDetails{
		AtExpires: time.Now().Add(time.Second * time.Duration(cfg.AccessTokenLiveTimeMinutes)).Unix(),
		RtExpires: time.Now().Add(time.Second - time.Duration(cfg.RefreshTokenLiveTimeDays)).Unix(),
	}

	atClaims := jwt.MapClaims{
		"exp":     td.AtExpires,
		"user_id": userID,
		"email":   email,
	}

	rtClaims := jwt.MapClaims{
		"exp":     td.RtExpires,
		"user_id": userID,
		"email":   email,
	}

	atWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	rtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	at, err := atWithClaims.SignedString([]byte(cfg.AccessTokenSecret))
	if err != nil {
		return nil, err
	}

	rt, err := rtWithClaims.SignedString([]byte(cfg.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	td.AccessToken = at
	td.RefreshToken = rt

	log.Println("token has been generated")

	return td, nil
}

func ValidateToken(signedToken string, cfg *config.ConfigToken) (*jwt.Token, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.AccessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, errors.New("expired token")
	}

	return token, nil
}

func RefreshToken(refreshToken string, cfg *config.ConfigToken) (*TokenDetails, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.RefreshTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userID := claims["user_id"].(string)
		email := claims["email"].(string)

		td, err := CreateToken(userID, email, cfg)
		if err != nil {
			return nil, err
		}

		log.Println("token has been refreshed")

		return td, nil
	} else {
		return nil, errors.New("refresh token expired")
	}
}
