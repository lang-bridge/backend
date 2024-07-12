package infra

import (
	"context"
	"log/slog"
	"os"

	slogotel "github.com/remychantenay/slog-otel"
	slogsentry "github.com/samber/slog-sentry/v2"
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

	if _, withSentry := os.LookupEnv("SENTRY_DSN"); withSentry {
		sentryHandler := slogsentry.Option{Level: slog.LevelWarn}.NewSentryHandler()
		handler = slogCombine{
			loggers: []slog.Handler{
				handler,
				sentryHandler,
			},
		}
	}

	handler = slogotel.New(handler, slogotel.WithNoTraceEvents(true))

	return slog.New(handler)
}

type slogCombine struct {
	loggers []slog.Handler
}

func (s slogCombine) Enabled(ctx context.Context, level slog.Level) bool {
	for _, logger := range s.loggers {
		if logger.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (s slogCombine) Handle(ctx context.Context, record slog.Record) error {
	for _, logger := range s.loggers {
		if !logger.Enabled(ctx, record.Level) {
			continue
		}
		if err := logger.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (s slogCombine) WithAttrs(attrs []slog.Attr) slog.Handler {
	loggers := make([]slog.Handler, len(s.loggers))
	for i, logger := range s.loggers {
		loggers[i] = logger.WithAttrs(attrs)
	}
	return slogCombine{loggers: loggers}
}

func (s slogCombine) WithGroup(name string) slog.Handler {
	loggers := make([]slog.Handler, len(s.loggers))
	for i, logger := range s.loggers {
		loggers[i] = logger.WithGroup(name)
	}
	return slogCombine{loggers: loggers}
}
