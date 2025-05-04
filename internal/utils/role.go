package roles

type Role string

const (
	RoleGuest   Role = "GUEST"
	RoleStaff   Role = "STAFF"
	RoleManager Role = "MANAGER"
	RoleAdmin   Role = "ADMIN"
)

var validRoles = map[Role]bool{
	RoleGuest:   true,
	RoleStaff:   true,
	RoleManager: true,
	RoleAdmin:   true,
}

func IsValidRole(role Role) bool {
	return validRoles[role]
}