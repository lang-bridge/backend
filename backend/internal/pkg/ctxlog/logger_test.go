package ctxlog

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONLogger(t *testing.T) {
	buffer := new(bytes.Buffer)

	logger := slog.New(slog.NewJSONHandler(buffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()
	ctx = WithLogger(ctx, logger)

	Debug(ctx, "Hello", slog.String("key", "value"))
	require.Contains(t, buffer.String(), `"level":"DEBUG","msg":"Hello","key":"value"`)

	buffer.Reset()
	Info(ctx, "Hello", slog.String("key2", "value2"))
	require.Contains(t, buffer.String(), `"level":"INFO","msg":"Hello","key2":"value2"`)

	buffer.Reset()
	Warn(ctx, "Hello", slog.String("key3", "value3"))
	require.Contains(t, buffer.String(), `"level":"WARN","msg":"Hello","key3":"value3"`)

	buffer.Reset()
	Error(ctx, "Hello", slog.String("key4", "value4"))
	require.Contains(t, buffer.String(), `"level":"ERROR","msg":"Hello","key4":"value4"`)
}
