package common

const (
	RoleJobseeker = "1"
	RoleRecruiter = "2"
	RoleAdmin     = "3"
)

func HasRole(role string, allowedRoles ...string) bool {
	if role == "" {
		return false
	}
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}
