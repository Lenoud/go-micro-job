package common

import (
	"strings"

	"user-service/user"
)

// Role constants — single source of truth for all modules.
const (
	RoleJobseeker = "1"
	RoleRecruiter = "2"
	RoleAdmin     = "3"
)

func HasRole(auth *user.UserContext, allowedRoles ...string) bool {
	if auth == nil {
		return false
	}
	for _, role := range allowedRoles {
		if auth.Role == role {
			return true
		}
	}
	return false
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
