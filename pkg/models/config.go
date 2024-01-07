package models

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerPort           int    `json:"serverPort" envconfig:"SERVER_PORT" required:"true"`
	RepositoriesRootPath string `json:"repositoriesRootPath" envconfig:"REPOSITORIES_ROOT_PATH" required:"true"`
	DbHost               string `json:"dbHost" envconfig:"DB_HOST" required:"true"`
	DbUser               string `json:"dbUser" envconfig:"DB_USER" required:"true"`
	DbPass               string `json:"dbPass" envconfig:"DB_PASS" required:"true"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from env, error: %w", err)
	}

	return cfg, nil
}
