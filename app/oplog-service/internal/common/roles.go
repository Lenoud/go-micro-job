package common

import (
	"oplog-service/oplog"

	sharedcommon "micro-shared/common"
)

const (
	RoleJobseeker = sharedcommon.RoleJobseeker
	RoleRecruiter = sharedcommon.RoleRecruiter
	RoleAdmin     = sharedcommon.RoleAdmin
)

func HasRole(auth *oplog.OpLogContext, allowedRoles ...string) bool {
	if auth == nil {
		return false
	}
	return sharedcommon.HasRole(auth.Role, allowedRoles...)
}
