package rbac

import (
	"github.com/mikespook/gorbac"
)

var Rbac *gorbac.RBAC

func init() {
	Rbac = gorbac.New()

	// roles
	rAdmin := gorbac.NewStdRole("admin")
	rUser := gorbac.NewStdRole("user")

	// Being on the box automatically enables permissions and access to
	// some endpoints
	rLocal := gorbac.NewStdRole("localhost")

	// regarding system settings
	prSettings := gorbac.NewStdPermission("r-settings") // read
	pwSettings := gorbac.NewStdPermission("w-settings") // write

	// regarding configuring other admin users
	prAdmins := gorbac.NewStdPermission("r-admins") // read
	pwAdmins := gorbac.NewStdPermission("w-admins") // write
	plAdmins := gorbac.NewStdPermission("l-admins") // list

	// regarding configuring other users
	prUsers := gorbac.NewStdPermission("r-users")
	pwUsers := gorbac.NewStdPermission("w-users")
	plUsers := gorbac.NewStdPermission("l-users")

	// regarding self
	pSelf := gorbac.NewStdPermission("self")

	// assigning permissions to roles
	rAdmin.Assign(prSettings)
	rAdmin.Assign(pwSettings)
	rAdmin.Assign(prAdmins)
	rAdmin.Assign(pwAdmins)
	rAdmin.Assign(plAdmins)
	rAdmin.Assign(prUsers)
	rAdmin.Assign(pwUsers)
	rAdmin.Assign(plUsers)
	rAdmin.Assign(pSelf)

	rUser.Assign(plUsers)
	rUser.Assign(pSelf)

	rLocal.Assign(plUsers)

	// add roles to rbac object
	Rbac.Add(rAdmin)
	Rbac.Add(rUser)
	rbac.Add(rLocalp)
}
