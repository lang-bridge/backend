package test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"platform/internal/app"
	"platform/internal/app/api"
	"platform/internal/infra"
	"platform/internal/repository/postgres"
)

func RunTest(t *testing.T, r interface{}) {
	require.NoError(t, infra.Init())
	cfg := app.Config{
		Logger: infra.LoggerConfig{
			Level:  slog.LevelDebug,
			Format: "console",
		},
		Database: postgres.DbConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "postgres",
			Database:        "langbridge",
			MaxOpenConns:    10,
			MaxIdleConns:    2,
			ConnMaxIdleTime: time.Minute,
		},
	}
	fxApp := fxtest.New(t, api.Module, fx.Replace(cfg), fx.Invoke(r))
	defer fxApp.RequireStop()
	fxApp.RequireStart()
}
