package config

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/dcwangmit01/goapi/app/util"
	uuid "github.com/satori/go.uuid"
)

type AppConfig struct {
	Settings *Settings `validate:"dive"`
	Users    []*User   `validate:"dive"`
}

func NewAppConfig() *AppConfig {

	// Only use if no current AppConfig exists.  Start from scratch and
	// initialize with a default admin login user/password.  Some fields
	// will be set to default, but most things will be set to zero values
	// and thus this AppConfig will not validate.  It will be up to the UI
	// to get field values from the first (admin) user, after which
	// Validate will be enforced

	adminPass, _ := bcrypt.GenerateFromPassword([]byte(DefaultAdminPassword), bcrypt.DefaultCost)

	return &AppConfig{
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

func AppConfigFromYaml(yamlString string) (*AppConfig, error) {
	ac := &AppConfig{}
	err := util.StructFromYamlStr(ac, yamlString)
	return ac, err
}

func (ac AppConfig) ToYaml() (string, error) {
	return util.StructToYamlStr(ac)
}

func (ac *AppConfig) GetUserByUsername(username string) (*User, error) {
	for _, u := range ac.Users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("Unable to find user by username")
}

func (ac *AppConfig) GetUserById(id string) (*User, error) {
	for _, u := range ac.Users {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("Unable to find user by id")
}

func (ac *AppConfig) AddUser(user *User) {
	ac.Users = append(ac.Users, user)
	return
}
