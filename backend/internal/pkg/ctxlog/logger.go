package ctxlog

import (
	"context"
	"github.com/samber/lo"
	"log/slog"
	"os"
)

type ctxMarker struct{}

var ctxMarkerKey = &ctxMarker{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	logger, _ := ctx.Value(ctxMarkerKey).(*slog.Logger)
	if logger == nil {
		if _, has := os.LookupEnv("PROD"); has {
			_, _ = os.Stderr.WriteString("logger is not set in context\n")
			return slog.Default()
		}
		panic("logger is not set in context")
	}
	return logger
}

func ErrorAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

func Debug(ctx context.Context, message string, args ...slog.Attr) {
	Log(ctx, slog.LevelDebug, message, args...)
}

func Info(ctx context.Context, message string, args ...slog.Attr) {
	Log(ctx, slog.LevelInfo, message, args...)
}

func Warn(ctx context.Context, message string, args ...slog.Attr) {
	Log(ctx, slog.LevelWarn, message, args...)
}

func Error(ctx context.Context, message string, args ...slog.Attr) {
	Log(ctx, slog.LevelError, message, args...)
}

func Log(ctx context.Context, level slog.Level, message string, args ...slog.Attr) {
	Logger(ctx).LogAttrs(ctx, level, message, args...)
}

func With(ctx context.Context, args ...slog.Attr) context.Context {
	logger := Logger(ctx).With(castArgs(args...)...)
	return WithLogger(ctx, logger)
}

func castArgs(args ...slog.Attr) []any {
	return lo.Map(args, func(item slog.Attr, _ int) any {
		return item
	})
}
