package domain

// UserRole represents a user's permission level.
type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleMember UserRole = "member"
)

// IsAdmin reports whether the role grants admin privileges.
func (r UserRole) IsAdmin() bool {
	return r == UserRoleAdmin
}
