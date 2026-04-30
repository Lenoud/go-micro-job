package common

import (
	"context"
	"testing"
)

func TestAuthFromContextBuildsUserContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "userId", float64(7))
	ctx = context.WithValue(ctx, "username", "admin")
	ctx = context.WithValue(ctx, "role", "3")

	auth, ok := AuthFromContext(ctx)
	if !ok {
		t.Fatal("AuthFromContext() ok = false, want true")
	}
	if auth.UserId != "7" || auth.Username != "admin" || auth.Role != "3" {
		t.Fatalf("AuthFromContext() = %+v, want userId=7 username=admin role=3", auth)
	}
}

func TestAuthFromContextRejectsMissingRole(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "userId", "7")
	ctx = context.WithValue(ctx, "username", "admin")

	if auth, ok := AuthFromContext(ctx); ok || auth != nil {
		t.Fatalf("AuthFromContext() = %+v, %v; want nil, false", auth, ok)
	}
}
