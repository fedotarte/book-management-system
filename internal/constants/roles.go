package constants

type Role string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
)

var Roles = struct {
	User      Role
	Moderator Role
	Admin     Role
}{
	User:      RoleUser,
	Moderator: RoleModerator,
	Admin:     RoleAdmin,
}
