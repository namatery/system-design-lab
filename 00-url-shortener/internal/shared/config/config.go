package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addr           string `env:"ADDR" envDefault:":8080"`
	PublicURL      string `evn:"PUBLIC_URL" envDefault:"http://localhost:8080"`
	ScyllaAddr     string `env:"SCYLLA_ADDR" envDefault:"localhost:9042"`
	ScyllaKeyspace string `env:"SCYLLA_KEYSPACE" envDefault:"url_shortener"`
}

func Load() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	return cfg, nil
}

func GetConfig() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}

	return cfg
}
