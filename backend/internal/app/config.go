package app

import (
	"fmt"
	"os"

	"go.uber.org/fx"
	"gopkg.in/yaml.v3"

	"platform/internal/api/http"
	"platform/internal/infra"
	"platform/internal/repository/postgres"
)

type Config struct {
	fx.Out

	Http     http.Config        `yaml:"http"`
	Logger   infra.LoggerConfig `yaml:"logger"`
	Database postgres.DbConfig  `yaml:"database"`
}

func ReadConfig() (Config, error) {
	var cfg Config
	path, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		path = "config.yaml"
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}
	err = yaml.Unmarshal([]byte(os.ExpandEnv(string(content))), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}
