package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var accessSecret = []byte("accessSecret")
var refreshSecret = []byte("refreshSecret")

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type JWTClaims struct {
	userID string
	email  string
	jwt.StandardClaims
}

func CreateToken(userID string, email string) (*TokenDetails, error) {
	td := &TokenDetails{
		AtExpires: time.Now().Add(time.Minute * time.Duration(60)).Unix(),
		RtExpires: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
	}

	atClaims := &JWTClaims{
		userID: userID,
		email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.AtExpires,
		},
	}

	rtClaims := &JWTClaims{
		userID: userID,
		email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.RtExpires,
		},
	}

	atWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	rtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	at, err := atWithClaims.SignedString(accessSecret)
	if err != nil {
		return nil, err
	}

	rt, err := rtWithClaims.SignedString(refreshSecret)
	if err != nil {
		return nil, err
	}

	td.AccessToken = at
	td.RefreshToken = rt

	return td, nil
}

func ValidateToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
