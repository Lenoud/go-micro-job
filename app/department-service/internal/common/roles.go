package common

import (
	"department-service/department"
	sharedcommon "micro-shared/common"
)

const (
	RoleJobseeker = sharedcommon.RoleJobseeker
	RoleRecruiter = sharedcommon.RoleRecruiter
	RoleAdmin     = sharedcommon.RoleAdmin
)

func HasRole(auth *department.DepartmentContext, allowedRoles ...string) bool {
	if auth == nil {
		return false
	}
	return sharedcommon.HasRole(auth.Role, allowedRoles...)
}
