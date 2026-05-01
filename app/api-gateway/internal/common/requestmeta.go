package common

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-Id"

type requestIDContextKey struct{}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil || strings.TrimSpace(requestID) == "" {
		return ctx
	}
	return context.WithValue(ctx, requestIDContextKey{}, strings.TrimSpace(requestID))
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(requestIDContextKey{}).(string); ok {
		return strings.TrimSpace(requestID)
	}
	return ""
}

func EnsureRequestID(r *http.Request) (string, *http.Request) {
	if r == nil {
		return "", r
	}

	requestID := strings.TrimSpace(r.Header.Get(requestIDHeader))
	if requestID == "" {
		requestID = uuid.NewString()
	}
	r.Header.Set(requestIDHeader, requestID)
	return requestID, r
}

func NewRequestMetaMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if next == nil {
			return nil
		}

		return func(w http.ResponseWriter, r *http.Request) {
			if r == nil {
				return
			}

			requestID, req := EnsureRequestID(r)
			if w != nil && requestID != "" {
				w.Header().Set(requestIDHeader, requestID)
			}
			next(w, req.WithContext(WithRequestID(req.Context(), requestID)))
		}
	}
}
