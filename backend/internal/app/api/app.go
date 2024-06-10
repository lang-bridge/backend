package api

import (
	"go.uber.org/fx"
	"platform/internal/api/http"
	"platform/internal/app"
	"platform/internal/infra"
)

var Module = fx.Module("api",
	fx.Provide(app.ReadConfig),
	infra.Module,
	http.Module,
)
