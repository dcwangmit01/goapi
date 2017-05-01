package rbac

import (
	"github.com/mikespook/gorbac"
)

var Rbac *gorbac.RBAC

var (
	// roles
	RoleAdmin = gorbac.NewStdRole("admin")
	RoleUser  = gorbac.NewStdRole("user")

	// Being on the box (localhost) automatically enables
	// permissions and access to some endpoints
	RoleLocal = gorbac.NewStdRole("local")

	// regarding system settings
	PermRSettings = gorbac.NewStdPermission("r-settings") // read
	PermWSettings = gorbac.NewStdPermission("w-settings") // write

	// regarding configuring other admin users
	PermRAdmins = gorbac.NewStdPermission("r-admins") // read
	PermWAdmins = gorbac.NewStdPermission("w-admins") // write
	PermLAdmins = gorbac.NewStdPermission("l-admins") // list

	// regarding configuring other users
	PermRUsers = gorbac.NewStdPermission("r-users")
	PermWUsers = gorbac.NewStdPermission("w-users")
	PermLUsers = gorbac.NewStdPermission("l-users")
)

func init() {
	Rbac = gorbac.New()

	// assigning permissions to roles
	RoleAdmin.Assign(PermRSettings)
	RoleAdmin.Assign(PermWSettings)
	RoleAdmin.Assign(PermRAdmins)
	RoleAdmin.Assign(PermWAdmins)
	RoleAdmin.Assign(PermLAdmins)
	RoleAdmin.Assign(PermRUsers)
	RoleAdmin.Assign(PermWUsers)
	RoleAdmin.Assign(PermLUsers)

	RoleUser.Assign(PermLUsers)

	RoleLocal.Assign(PermLUsers)

	// add roles to rbac object
	Rbac.Add(RoleAdmin)
	Rbac.Add(RoleUser)
	Rbac.Add(RoleLocal)
}
