package common

import (
	"strings"

	sharedcommon "micro-shared/common"
	"user-service/user"
)

// Role constants — single source of truth for all modules.
const (
	RoleJobseeker = sharedcommon.RoleJobseeker
	RoleRecruiter = sharedcommon.RoleRecruiter
	RoleAdmin     = sharedcommon.RoleAdmin
)

func HasRole(auth *user.UserContext, allowedRoles ...string) bool {
	if auth == nil {
		return false
	}
	return sharedcommon.HasRole(auth.Role, allowedRoles...)
}

func IsAdminContext(auth *user.UserContext) bool {
	return HasRole(auth, RoleAdmin)
}

func DetailTargetUserID(auth *user.UserContext, requestedUserID string) (string, bool) {
	if auth == nil || strings.TrimSpace(auth.UserId) == "" {
		return "", false
	}
	trimmedUserID := strings.TrimSpace(requestedUserID)
	if IsAdminContext(auth) && trimmedUserID != "" {
		return trimmedUserID, true
	}
	return auth.UserId, true
}

func ScopedMutationUserID(auth *user.UserContext, requestedUserID string) (string, bool) {
	if auth == nil || strings.TrimSpace(auth.UserId) == "" {
		return "", false
	}
	trimmedUserID := strings.TrimSpace(requestedUserID)
	if IsAdminContext(auth) {
		if trimmedUserID != "" {
			return trimmedUserID, true
		}
		return auth.UserId, true
	}
	if trimmedUserID != "" && trimmedUserID != auth.UserId {
		return "", false
	}
	return auth.UserId, true
}
