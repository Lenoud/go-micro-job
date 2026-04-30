package common

import "department-service/department"

const (
	RoleJobseeker = "1"
	RoleRecruiter = "2"
	RoleAdmin     = "3"
)

func HasRole(auth *department.DepartmentContext, allowedRoles ...string) bool {
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
