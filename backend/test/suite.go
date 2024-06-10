package test

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"log/slog"
	"math/rand/v2"
	"platform/internal/app"
	"platform/internal/app/api"
	"platform/internal/infra"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/service/keys"
	"platform/pkg/ctxlog"
	"testing"
)

func RunTest(t *testing.T, r interface{}) {
	cfg := app.Config{
		Logger: infra.LoggerConfig{
			Level:  slog.LevelDebug,
			Format: "console",
		},
	}
	fxApp := fxtest.New(t, api.Module, fx.Provide(newStub), fx.Replace(cfg), fx.Invoke(r))
	defer fxApp.RequireStop()
	fxApp.RequireStart()
}

var _ keys.Service = (*stubService)(nil)

type stubService struct {
}

func newStub() keys.Service {
	return stubService{}
}

func (s stubService) CreateKey(ctx context.Context, _ keys.CreateKeyParam) (keys.KeyView, error) {
	ctxlog.Info(ctx, "test info message")
	return keys.KeyView{
		Key: key.Key{
			ID: key.ID(rand.Int64()),
		},
	}, nil
}
