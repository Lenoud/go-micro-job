package common

import (
	"context"

	departmentclient "department-service/departmentClient"
)

func DepartmentAuthFromContext(ctx context.Context) (*departmentclient.DepartmentContext, bool) {
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
	return &departmentclient.DepartmentContext{
		UserId:   userID,
		Username: username,
		Role:     role,
	}, true
}
