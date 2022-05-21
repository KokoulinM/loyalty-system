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
			panic(err)
		}

		fmt.Println(jwt.Raw)

		time.Sleep(30 * time.Second)
	}
}
