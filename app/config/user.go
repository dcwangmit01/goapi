package config

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dcwangmit01/goapi/app/jwt"
	rbac "github.com/dcwangmit01/goapi/app/rbac"
	gorbac "github.com/mikespook/gorbac"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id string `validate:"uuid4"`
	// Username is globally unique from the user's perspective.  It must
	// either be in email address format, or the special case 'admin'
	Username     string `validate:"required,email|eq=admin"`
	Name         string `validate:"required,printascii"` // TODO: plan on escaping and handling all printable ASCII
	PasswordHash string `validate:"required"`            // Hashed and Salted by the bcrypt library
	Role         string `validate:"eq=user|eq=admin"`
	Phone        string `validate:"phone,min=7"`
}

func NewUser() *User {
	return &User{
		Id:   uuid.NewV4().String(),
		Role: "user",
	}
}

func (u *User) GenerateJwt(secondsToExpiration int64) (string, error) {
	return jwt.CreateJwtWithIdRole(u.Id, u.Role, secondsToExpiration)
}

func (u *User) GetRole() (gorbac.Role, error) {
	role, _, err := rbac.Rbac.Get(u.Role)
	return role, err
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		u.PasswordHash = string(hash)
	}
	return err
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

const userContextKey = "user"

func UserNewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}
