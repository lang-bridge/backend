package httputil

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"platform/pkg/ctxlog"
	"platform/pkg/httputil/httperr"
)

func WrapError(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if httperr.IsBadRequest(err) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ctxlog.Error(r.Context(), "request finished with unhandled error", ctxlog.ErrorAttr(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func WithLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctxlog.WithLogger(r.Context(), logger))
			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				ctxlog.Error(r.Context(), "recovered from panic",
					slog.Any("panic", rvr),
					slog.String("stack", string(debug.Stack())),
				)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
