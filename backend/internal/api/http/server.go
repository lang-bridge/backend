package http

import (
	"context"
	"errors"
	"fmt"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	slogchi "github.com/samber/slog-chi"
	"log/slog"
	"net/http"
	"platform/pkg/ctxlog"
	"platform/pkg/httputil"
)

type Server struct {
	*http.Server
}

type Registerer interface {
	Register(router chi.Router)
}

func (s *Server) Start() error {
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("couldn't start http server: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("couldn't gracefully stop http server: %w", err)
	}
	ctxlog.Debug(ctx, "http server stopped")
	return nil
}

type Config struct {
	Port int `yaml:"port"`
}

func NewServer(cfg Config, router chi.Router) (*Server, error) {
	if cfg.Port <= 0 {
		return nil, errors.New("http port cannot be less 0")
	}

	return &Server{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: router,
		},
	}, nil
}

func NewRouter(logger *slog.Logger, registrars ...Registerer) chi.Router {
	r := chi.NewRouter()

	r.Use(httputil.WithLogger(logger))

	r.Use(otelchi.Middleware("api", otelchi.WithChiRoutes(r)))

	r.Use(sentryhttp.New(sentryhttp.Options{}).Handle)

	r.Use(slogchi.New(logger))

	r.Use(httputil.Recovery)

	for _, registrar := range registrars {
		registrar.Register(r)
	}
	return r
}
