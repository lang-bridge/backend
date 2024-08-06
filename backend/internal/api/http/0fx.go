package http

import (
	"context"
	"log/slog"
	"platform/internal/api/http/hello"
	"platform/internal/pkg/ctxlog"

	"go.uber.org/fx"

	"platform/internal/api/http/keys"
)

var Module = fx.Module("http",
	ServerModule,
	RoutersModule,
)

var Invoke = fx.Invoke(
	RunServer,
)

var ServerModule = fx.Provide(
	fx.Annotate(NewRouter, fx.ParamTags("", groupTag)),
	NewServer,
)

const groupTag = `group:"http_handlers"`

var RoutersModule = fx.Provide(
	fx.Annotate(keys.NewRouter, fx.As(new(Registerer)), fx.ResultTags(groupTag)),
	fx.Annotate(hello.NewRouter, fx.As(new(Registerer)), fx.ResultTags(groupTag)),
)

// RunServer adds to lifecycle a new fx.Hook
// which runs and stops http srv
func RunServer(lc fx.Lifecycle, srv *Server, logger *slog.Logger) {
	lc.Append(
		// http server hook
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := srv.Start(); err != nil {
						logger.ErrorContext(ctx, "couldn't start http server", ctxlog.ErrorAttr(err))
					}
				}()
				return nil
			},
			OnStop: srv.Shutdown,
		},
	)
}
