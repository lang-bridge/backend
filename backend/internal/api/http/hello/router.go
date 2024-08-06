package hello

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"platform/pkg/ctxlog"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (s *Router) Register(router chi.Router) {
	router.Get("/hello", s.Hello)
}

func (s *Router) Hello(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		for _, value := range values {
			ctxlog.Info(r.Context(), "Header", slog.String("key", key), slog.String("value", value))
		}
	}

	_, _ = w.Write([]byte("Hello, world!"))
	w.WriteHeader(http.StatusOK)
}
