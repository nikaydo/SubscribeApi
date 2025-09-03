package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	Port       string `env:"PORT"`
	Host       string `env:"HOST"`
	Postgresql string `env:"POSTGRESQL"`
}

func ReadEnv() (Env, error) {
	err := godotenv.Load()
	var cfg Env
	if err != nil {
		return cfg, fmt.Errorf("error read env: %w", err)
	}
	if err := env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("error read env: %w", err)
	}
	return cfg, nil
}
