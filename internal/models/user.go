package models

const (
	AdultAge uint8 = 18
)

type UserRole int8

const (
	RoleUser UserRole = iota
	RoleAdmin
)

const (
	userStr  string = "user"
	adminStr string = "admin"
)

var UserRolesStr = map[string]UserRole{
	userStr:  RoleUser,
	adminStr: RoleAdmin,
}

func (r UserRole) String() string {
	switch r {
	case RoleUser:
		return userStr

	case RoleAdmin:
		return adminStr

	default:
		return "undefined"
	}
}

func (r UserRole) FromString(role string) UserRole {
	return UserRolesStr[role]
}

type User struct {
	ID           string   `db:"id"`
	Username     string   `db:"username"`
	Name         string   `db:"name"`
	Role         UserRole `db:"role"`
	Age          uint8    `db:"age"`
	PasswordHash string   `db:"password_hash"`
}
