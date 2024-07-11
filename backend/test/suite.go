package test

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"log/slog"
	"platform/internal/app"
	"platform/internal/app/api"
	"platform/internal/infra"
	"testing"
	"time"
)

func RunTest(t *testing.T, r interface{}) {
	cfg := app.Config{
		Logger: infra.LoggerConfig{
			Level:  slog.LevelDebug,
			Format: "console",
		},
		Database: infra.DbConfig{
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
