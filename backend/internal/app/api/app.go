package api

import (
	"go.uber.org/fx"
	"platform/internal/api/http"
	"platform/internal/app"
	"platform/internal/infra"
	"platform/internal/repository/postgres"
	"platform/internal/translations/service/keys"
)

var Module = fx.Module("api",
	fx.Provide(app.ReadConfig),
	infra.Module,
	http.Module,
	postgres.Module,

	// domain
	keys.Module,
)
