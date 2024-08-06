package hello

import (
	"github.com/go-chi/chi/v5"
	"net/http"
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
	_, _ = w.Write([]byte("Hello, world!"))
	w.WriteHeader(http.StatusOK)
}
