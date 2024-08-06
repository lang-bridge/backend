package httputil

import (
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"platform/internal/pkg/ctxlog"
	"platform/internal/pkg/metactx"
	"platform/internal/pkg/types"
)

const xUserHeader = "X-User"

func WithMetaCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Header[xUserHeader]
		if len(value) > 0 {
			r = enrichUserID(r, value[0])
		}

		next.ServeHTTP(w, r)
	})
}

func enrichUserID(r *http.Request, value string) *http.Request {
	userID, err := uuid.Parse(value)
	if err != nil {
		ctxlog.Error(r.Context(), "can't parse userID", ctxlog.ErrorAttr(err))
		return r
	}
	ctx := r.Context()
	ctx = metactx.WithUserID(ctx, types.UserID(userID))
	ctx = ctxlog.With(ctx, slog.String("user-id", userID.String()))

	return r.WithContext(ctx)
}
