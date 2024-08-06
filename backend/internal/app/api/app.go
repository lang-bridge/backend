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

var Invoke = fx.Options(
	http.Invoke,
)

// App is a set of invokes and modules
var App = fx.Options(
	Module,
	Invoke,
)

// New returns a new instance of the App
//
//	@title			LangBridge API
//	@version		1.0
//	@contact.name	LangBridge Support
//	@contact.url	http://langbridge.io/support
//	@contact.email	support@langbridge.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			https://api.langbridge.io
func New() *fx.App {
	return fx.New(App)
}

// Run runs the App
func Run() {
	err := infra.Init()
	if err != nil {
		panic(err)
	}
	New().Run()
}
