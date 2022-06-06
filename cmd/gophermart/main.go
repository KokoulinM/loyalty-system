package main

import (
	"fmt"
	"time"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/auth"
)

func main() {
	token, err := auth.CreateToken("1", "")
	if err != nil {
		panic(err)
	}

	for {
		jwt, err := auth.ValidateToken(token.AccessToken)
		if err != nil {
			token, err = auth.RefreshToken(token.RefreshToken)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(jwt.Raw)
		}

		time.Sleep(1 * time.Second)
	}
}
