package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	ServerAddress              = ":8080"
	DataBaseURI                = "postgres://postgres:postgres@localhost:5432/gophermartdb?sslmode=disable"
	AccrualSystemAddress       = ""
	AccessTokenSecret          = ""
	RefreshTokenSecret         = ""
	AccessTokenLiveTimeMinutes = 60
	RefreshTokenLiveTimeDays   = 7
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	DataBase             ConfigDatabase
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Token                ConfigToken
}

type ConfigToken struct {
	AccessTokenSecret          string `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret         string `env:"REFRESH_TOKEN_SECRET"`
	AccessTokenLiveTimeMinutes int    `env:"ACCESS_TOKEN_LIVE_TIME_MINUTES"`
	RefreshTokenLiveTimeDays   int    `env:"REFRESH_TOKEN_LIVE_TIME_DAYS"`
}

type ConfigDatabase struct {
	DataBaseURI string `env:"DATABASE_URI"`
}

func New() *Config {
	dbCfg := ConfigDatabase{
		DataBaseURI: DataBaseURI,
	}

	tokenConfig := NewConfigToken()

	cfg := Config{
		ServerAddress:        ServerAddress,
		DataBase:             dbCfg,
		AccrualSystemAddress: AccrualSystemAddress,
	}

	err := env.Parse(&cfg.DataBase)
	if err != nil {
		log.Fatal(err)
	}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("an error occurred during parsing Config: %s", err)
	}

	if checkExists("a") {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "ServerAddress")
	}

	if checkExists("d") {
		flag.StringVar(&cfg.DataBase.DataBaseURI, "d", cfg.DataBase.DataBaseURI, "DataBaseURI")
	}

	if checkExists("r") {
		flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "AccrualSystemAddress")
	}

	cfg.Token = tokenConfig

	flag.Parse()

	return &cfg
}

func NewConfigToken() ConfigToken {
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
