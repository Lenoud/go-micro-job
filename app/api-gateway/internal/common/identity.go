package common

import (
	"context"
	"strconv"
	"strings"

	userclient "user-service/userClient"
)

func AuthFromContext(ctx context.Context) (*userclient.UserContext, bool) {
	userID, ok := claimString(ctx, "userId")
	if !ok {
		return nil, false
	}
	username, ok := claimString(ctx, "username")
	if !ok {
		return nil, false
	}
	role, ok := claimString(ctx, "role")
	if !ok {
		return nil, false
	}
	return &userclient.UserContext{
		UserId:   userID,
		Username: username,
		Role:     role,
	}, true
}

func claimString(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}

	var str string
	switch v := ctx.Value(key).(type) {
	case string:
		str = v
	case float64:
		str = strconv.FormatInt(int64(v), 10)
	case int:
		str = strconv.Itoa(v)
	case int64:
		str = strconv.FormatInt(v, 10)
	default:
		return "", false
	}

	str = strings.TrimSpace(str)
	if str == "" {
		return "", false
	}
	return str, true
}
