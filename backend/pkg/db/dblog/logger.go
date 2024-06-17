package dblog

import (
	"context"
	"github.com/jackc/pgx/v5/tracelog"
	"log/slog"
	"platform/pkg/ctxlog"
	"sort"
)

var _ tracelog.Logger = (*slogLogger)(nil)

func NewLogger() tracelog.Logger {
	return &slogLogger{}
}

type slogLogger struct {
}

func (s *slogLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var attrs []slog.Attr
	for _, k := range keys {
		attrs = append(attrs, slog.Any(k, data[k]))
	}
	ctxlog.Log(ctx, translateLevel(level), msg, attrs...)
}

func translateLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	case tracelog.LogLevelNone:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
