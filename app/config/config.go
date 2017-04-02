package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
)

const (
	defaultAdminUser = "admin"
	defaultAdminPass = "password"
)

type AppConfig struct {
	Settings *Settings `validate:"dive"`
	Users    []*User   `validate:"dive"`
}

type Settings struct {
	Initialized bool
	Debug       bool
	LogLevel    string `validate:"eq=DEBUG|eq=INFO"`
}
type User struct {
	Id           string `validate:"uuid4"`
	Email        string `validate:"required,email"`
	Name         string `validate:"required,printascii"` // TODO: plan on escaping and handling all printable ASCII
	PasswordHash string `validate:"required,base64"`     // Hashed and Salted by the bcrypt library
	Role         string `validate:"eq=USER|eq=ADMIN"`
	PhoneNumber  string `validate:"phone,min=7"`
}

func (u *User) ValidatePassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewUser() *User {
	return &User{
		Id:   uuid.NewV4().String(),
		Role: "USER",
	}
}

func NewAppConfig() *AppConfig {

	// If no current AppConfig exists, start from scratch and initialize
	// with a default admin login user/password.  Some fields will be set
	// to default, but most things will be set to zero values and thus this
	// AppConfig will not validate.  It will be up to the UI to get field
	// values from the first (admin) user, after which Validate will be
	// enforced

	adminPass, _ := bcrypt.GenerateFromPassword([]byte(defaultAdminPass), bcrypt.DefaultCost)

	return &AppConfig{
		&Settings{
			LogLevel: "DEBUG",
		},
		[]*User{
			&User{ // default admin user
				Id:           uuid.NewV4().String(),
				Role:         "ADMIN",
				PasswordHash: string(adminPass),
			},
		},
	}
}

func ParseAppConfig(yamlString string) (*AppConfig, error) {
	// Parse appConfig from a string
	ac := &AppConfig{}

	err := yaml.Unmarshal([]byte(yamlString), ac)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return ac, err

}

func (ac *AppConfig) Dump() (string, error) {
	d, err := yaml.Marshal(ac)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d), err
}

func main() {

	ac := NewAppConfig()
	ac.Users = append(ac.Users, NewUser())

	ac.Users[0].Email = "user1@gmail.com"
	ac.Users[1].Email = "asdf"

	validate := validator.New()

	// A phone number contains no dashes, spaces, and parens. Valid keys
	// all keypad [0-9#*] as well as "+" to designate country code, and "w"
	// to signify a .5 second wait.
	var phoneRegexString = "^+?[0-9#*w]+$"
	var phoneRegex = regexp.MustCompile(phoneRegexString)
	var isPhone = func(fl validator.FieldLevel) bool {
		return phoneRegex.MatchString(fl.Field().String())
	}
	validate.RegisterValidation("phone", isPhone)

	err := validate.Struct(ac)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		fmt.Println(errs)
	}

	ac1Str, err := ac.Dump()
	ac2, err := ParseAppConfig(ac1Str)
	ac2Str, err := ac2.Dump()
	fmt.Println(ac1Str)
	fmt.Println(ac2Str)
}
