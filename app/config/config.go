package config

import (
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
)

var validate = validator.New()

func init() {
	// Register an additional validator called "phone"
	// A phone number contains no dashes, spaces, and parens. Valid keys
	// all keypad [0-9#*] as well as "+" to designate country code, and "w"
	// to signify a .5 second wait.
	var phoneRegexString = "^+?[0-9#*w]+$"
	var phoneRegex = regexp.MustCompile(phoneRegexString)
	var isPhone = func(fl validator.FieldLevel) bool {
		return phoneRegex.MatchString(fl.Field().String())
	}

	validate.RegisterValidation("phone", isPhone)
}

const (
	DefaultAdminUser = "admin"
	DefaultAdminPass = "password"
)

type AppConfig struct {
	Settings *Settings `validate:"dive"`
	Users    []*User   `validate:"dive"`
}

type Settings struct {
	Initialized bool
	Debug       bool
	LogLevel    string `validate:"eq=DEBUG|eq=INFO|eq=WARNING|eq=ERROR|eq=FATAL|eq=PANIC"`
}

type User struct {
	Id           string `validate:"uuid4"`
	Email        string `validate:"required,email"`
	Name         string `validate:"required,printascii"` // TODO: plan on escaping and handling all printable ASCII
	PasswordHash string `validate:"required"`            // Hashed and Salted by the bcrypt library
	Role         string `validate:"eq=USER|eq=ADMIN"`
	Phone        string `validate:"phone,min=7"`
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

func NewUser() *User {
	return &User{
		Id:   uuid.NewV4().String(),
		Role: "USER",
	}
}

func NewAppConfig() *AppConfig {

	// Only use if no current AppConfig exists.  Start from scratch and
	// initialize with a default admin login user/password.  Some fields
	// will be set to default, but most things will be set to zero values
	// and thus this AppConfig will not validate.  It will be up to the UI
	// to get field values from the first (admin) user, after which
	// Validate will be enforced

	adminPass, _ := bcrypt.GenerateFromPassword([]byte(DefaultAdminPass), bcrypt.DefaultCost)

	return &AppConfig{
		&Settings{
			LogLevel: "DEBUG",
		},
		[]*User{
			&User{ // default admin user
				Id:           uuid.NewV4().String(),
				Role:         "ADMIN",
				Name:         DefaultAdminUser,
				PasswordHash: string(adminPass),
			},
		},
	}
}

func ParseAppConfig(yamlString string) (*AppConfig, error) {
	// Parse from a string and populate a new nested AppConfig struct
	ac := &AppConfig{}

	err := yaml.Unmarshal([]byte(yamlString), ac)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}
	return ac, err

}

func (ac *AppConfig) AddUser(user *User) {
	ac.Users = append(ac.Users, user)
	return
}

func (ac *AppConfig) Dump() (string, error) {
	d, err := yaml.Marshal(ac)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d), err
}

func ValidateStruct(s interface{}) map[string]string {
	// returns nil if no errors
	//   or a map with key=namespace, value=validation_tag_that_failed

	err := validate.Struct(s)
	if err != nil {
		fieldErrors := err.(validator.ValidationErrors)

		// translate the fieldErrors above from odd types to
		// map[string]string
		m := make(map[string]string)
		for i := 0; i < len(fieldErrors); i++ {
			fe := fieldErrors[i]
			m[fe.Namespace()] = fe.Tag()
		}
		return m
	}
	return nil
}
