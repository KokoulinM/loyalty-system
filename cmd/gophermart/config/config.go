package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	ServerAddress              = ":8080"
	DataBaseURI                = "postgres://postgres:postgres@localhost:5432/gophermartdb?sslmode=disable"
	AccrualSystemAddress       = "http://localhost:8080"
	AccessTokenSecret          = ""
	RefreshTokenSecret         = ""
	AccessTokenLiveTimeMinutes = 60
	RefreshTokenLiveTimeDays   = 7
	NumOfWorkers               = 10
	PoolBuffer                 = 1000
	MaxJobRetryCount           = 5
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DataBase             ConfigDatabase
	Token                ConfigToken
	WorkerPool           ConfigWorkerPool
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

type ConfigWorkerPool struct {
	NumOfWorkers     int `env:"num_of_workers"`
	PoolBuffer       int `env:"pool_buffer"`
	MaxJobRetryCount int `env:"max_job_retry_count"`
}

func New() *Config {
	dbCfg := ConfigDatabase{
		DataBaseURI: DataBaseURI,
	}

	tokenConfig := NewConfigToken()

	wpConf := ConfigWorkerPool{
		NumOfWorkers:     NumOfWorkers,
		PoolBuffer:       PoolBuffer,
		MaxJobRetryCount: MaxJobRetryCount,
	}

	cfg := Config{
		ServerAddress:        ServerAddress,
		DataBase:             dbCfg,
		AccrualSystemAddress: AccrualSystemAddress,
		WorkerPool:           wpConf,
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
