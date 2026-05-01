package common

import (
	"context"

	oplogclient "oplog-service/oplogclient"
)

func OpLogAuthFromContext(ctx context.Context) (*oplogclient.OpLogContext, bool) {
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
	return &oplogclient.OpLogContext{
		UserId:   userID,
		Username: username,
		Role:     role,
	}, true
}
