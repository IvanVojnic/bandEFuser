// Package config to db
package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

// Config struct used to declare db connection
type Config struct {
	USER             string `env:"USER" envDefault:"postgres"`
	PASSWORD         string `env:"PASSWORD" envDefault:"postgres"`
	PostgresPort     string `env:"POSTGRES_PORT,notEmpty" envDefault:"5432"`
	PostgresHost     string `env:"POSTGRES_HOST,notEmpty" envDefault:"localhost"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,notEmpty" envDefault:"postgres"`
	PostgresUser     string `env:"POSTGRES_USER,notEmpty" envDefault:"postgres"`
	PostgresDB       string `env:"POSTGRES_DB,notEmpty" envDefault:"postgres"`

	Port string `env:"PORT" envDefault:"8000"`
	Host string `env:"HOST" envDefault:"0.0.0.0"`
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
