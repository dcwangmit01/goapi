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
	defaultAdminUuid = "11111111-1111-1111-1111-111111111111" // 128 bits set to 1
)

type AppConfig struct {
	Settings *Settings `validate:"dive"`
	Users    []*User   `validate:"dive"`
}

type Settings struct {
	Debug    bool
	LogLevel string `validate:"eq=DEBUG|eq=INFO"`
}
type User struct {
	Id             uuid.UUID `validate:"uuid"`
	Email          string    `validate:"required,email"`
	Name           string    `validate:"required,printascii"` // TODO: plan on escaping and handling all printable ASCII
	HashedPassword string    `validate:"required,base64"`     // Hashed and Salted by the bcrypt library
	Role           string    `validate:"eq=user|eq=admin"`
	PhoneNumber    string    `validate:"phone,min=7"`
}

func (u *User) ValidatePassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewUser() *User {
	return &User{
		Id:   uuid.NewV4(),
		Role: "USER",
	}
}

func NewAppConfig() *AppConfig {

	// If starting with an empty AppConfig, initialize with a default admin
	// login user/password.  Some fields will be set to default, but most
	// things will be set to zero values and thus this AppConfig will not
	// validate.  It will be up to the UI to get field values from the
	// first (admin) user, after which Validate will be enforced

	adminUuid, _ := uuid.FromString(defaultAdminUuid)
	adminPass, _ := bcrypt.GenerateFromPassword([]byte(defaultAdminPass), bcrypt.DefaultCost)

	return &AppConfig{
		&Settings{
			LogLevel: "DEBUG",
		},
		[]*User{
			&User{ // default admin user
				Id:             adminUuid,
				Role:           "ADMIN",
				HashedPassword: string(adminPass),
			},
		},
	}
}

func main() {
	// t := T{}

	// err := yaml.Unmarshal([]byte(data), &t)
	// if err != nil {
	//         log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t:\n%v\n\n", t)

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

	d, err := yaml.Marshal(&ac)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

}
