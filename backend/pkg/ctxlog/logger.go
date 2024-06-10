package ctxlog

import (
	"context"
	"log/slog"
)

type ctxMarker struct{}

var ctxMarkerKey = &ctxMarker{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	logger, _ := ctx.Value(ctxMarkerKey).(*slog.Logger)
	return logger
}

func ErrorAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

func Debug(ctx context.Context, message string, args ...slog.Attr) {
	Logger(ctx).DebugContext(ctx, message, args)
}

func Info(ctx context.Context, message string, args ...slog.Attr) {
	Logger(ctx).InfoContext(ctx, message, args)
}

func Warn(ctx context.Context, message string, args ...slog.Attr) {
	Logger(ctx).WarnContext(ctx, message, args)
}

func Error(ctx context.Context, message string, args ...slog.Attr) {
	Logger(ctx).ErrorContext(ctx, message, args)
}
