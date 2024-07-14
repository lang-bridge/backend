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

var Invoke = fx.Invoke(
	http.RunServer,
)

// App is a set of invokes and modules
var App = fx.Options(
	Module,
	Invoke,
)

// New returns a new instance of the App
//
//	@title						Swagger Example API
//	@version					1.0
//	@description				This is a sample server celler server.
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@securityDefinitions.basic	BasicAuth
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
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
