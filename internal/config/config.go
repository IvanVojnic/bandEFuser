// Package config to db
package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

// Config struct used to declare db connection
type Config struct {
	USER        string `env:"USER" envDefault:"postgres"`
	PostgresURL string `env:"pUrl" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	PASSWORD    string `env:"PASSWORD" envDefault:"postgres"`
	PORT        int    `env:"PORT" envDefault:"5432"`
	DB          string `env:"DB" envDefault:"postgres"`
}

// AuthConfig struct used to declare auth var
/*type AuthConfig struct {
	SigningKey      string `env:"SigningKey" envDefault:"barband"`
	TokenRTDuration int    `env:"TokenRTDuration" envDefault:"3600000000000000"`
	TokenATDuration int    `env:"TokenATDuration" envDefault:"3600000000000"`
}*/

// NewConfig used to init config to db
func NewConfig() (*Config, error) {
	Cfg := &Config{}
	if err := env.Parse(Cfg); err != nil {
		return nil, fmt.Errorf("config - NewConfig: %v", err)
	}
	return Cfg, nil
}
