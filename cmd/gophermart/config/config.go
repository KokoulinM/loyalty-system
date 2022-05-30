package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	ServerAddress              = "localhost:8080"
	DataBaseURI                = "user=postgres password=postgres sslmode=disable"
	AccrualSystemAddress       = ""
	AccessTokenSecret          = ""
	RefreshTokenSecret         = ""
	AccessTokenLiveTimeMinutes = 60
	RefreshTokenLiveTimeDays   = 7
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Token                ConfigToken
}

type ConfigToken struct {
	AccessTokenSecret          string `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret         string `env:"REFRESH_TOKEN_SECRET"`
	AccessTokenLiveTimeMinutes int    `env:"ACCESS_TOKEN_LIVE_TIME_MINUTES"`
	RefreshTokenLiveTimeDays   int    `env:"REFRESH_TOKEN_LIVE_TIME_DAYS"`
}

func New() Config {
	cfg := Config{
		ServerAddress:        ServerAddress,
		DataBaseURI:          DataBaseURI,
		AccrualSystemAddress: AccrualSystemAddress,
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("an error occurred during parsing Config: %s", err)
	}

	if checkExists("a") {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "ServerAddress")
	}

	if checkExists("d") {
		flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "DataBaseURI")
	}

	if checkExists("r") {
		flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "AccrualSystemAddress")
	}

	tokenConfig := newConfigToken()

	cfg.Token = tokenConfig

	flag.Parse()

	return cfg
}

func newConfigToken() ConfigToken {
	cfg := ConfigToken{
		AccessTokenSecret:          AccessTokenSecret,
		RefreshTokenSecret:         RefreshTokenSecret,
		AccessTokenLiveTimeMinutes: AccessTokenLiveTimeMinutes,
		RefreshTokenLiveTimeDays:   RefreshTokenLiveTimeDays,
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("an error occurred during parsing ConfigToken: %s", err)
	}

	return cfg
}

func checkExists(f string) bool {
	return flag.Lookup(f) == nil
}
