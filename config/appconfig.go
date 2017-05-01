package config

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/dcwangmit01/goapi/util"
	uuid "github.com/satori/go.uuid"
)

type appConfig struct {
	Settings *Settings `validate:"dive"`
	Users    []*User   `validate:"dive"`
}

func NewAppConfig() *appConfig {

	// Only use if no current appConfig exists.  Start from scratch and
	// initialize with a default admin login user/password.  Some fields
	// will be set to default, but most things will be set to zero values
	// and thus this appConfig will not validate.  It will be up to the UI
	// to get field values from the first (admin) user, after which
	// Validate will be enforced

	adminPass, _ := bcrypt.GenerateFromPassword([]byte(DefaultAdminPassword), bcrypt.DefaultCost)

	return &appConfig{
		&Settings{
			LogLevel: "debug",
		},
		[]*User{
			&User{ // default admin user
				Id:           uuid.NewV4().String(),
				Username:     "admin",
				Role:         "admin",
				Name:         DefaultAdminUsername,
				PasswordHash: string(adminPass),
			},
		},
	}
}

func AppConfigFromYaml(yamlString string) (*appConfig, error) {
	ac := &appConfig{}
	err := util.StructFromYamlStr(ac, yamlString)
	return ac, err
}

func (ac appConfig) ToYaml() (string, error) {
	return util.StructToYamlStr(ac)
}

func (ac *appConfig) GetUserByUsername(username string) (*User, error) {
	for _, u := range ac.Users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("Unable to find user by username")
}

func (ac *appConfig) GetUserById(id string) (*User, error) {
	for _, u := range ac.Users {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("Unable to find user by id")
}

func (ac *appConfig) AddUser(user *User) {
	ac.Users = append(ac.Users, user)
	return
}
