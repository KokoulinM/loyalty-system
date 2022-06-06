package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var accessTokenSecret = []byte("accessTokenSecret")
var refreshTokenSecret = []byte("refreshTokenSecret")

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID string, email string) (*TokenDetails, error) {
	td := &TokenDetails{
		AtExpires: time.Now().Add(time.Minute * time.Duration(1)).Unix(),
		RtExpires: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
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

	at, err := atWithClaims.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return nil, err
	}

	rt, err := rtWithClaims.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return nil, err
	}

	td.AccessToken = at
	td.RefreshToken = rt

	return td, nil
}

func ValidateToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	err = token.Claims.Valid()
	if err != nil {
		return nil, err
	}

	return token, nil
}
