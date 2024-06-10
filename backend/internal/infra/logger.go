package infra

import (
	"log/slog"
	"os"
)

type LoggerConfig struct {
	Level  slog.Level `yaml:"level"`
	Format string     `yaml:"format"`
}

func NewLogger(config LoggerConfig) *slog.Logger {
	var handler slog.Handler
	switch config.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: config.Level,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: config.Level,
		})
	}
	return slog.New(handler)
}
