package main

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
)

type AppConfig struct {
	Config struct {
		Debug    bool
		LogLevel string `validate:"eq=DEBUG|eq=INFO"`
	}
	Users  []User `validate:"dive"`
}



// Unique ID is email
type User struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,printascii"` // TODO: plan on escaping and handling all printable ASCII
	Password struct {
		// http://stackoverflow.com/questions/23039458/equivalent-salt-and-hash-in-golang
		Salt string `validate:"required,base64"`
		Hash string `validate:"required,base64"`
	}
	Role        string `validate:"eq=user|eq=admin"`
	PhoneNumber string `validate:"phone,min=7"`
}

func main() {
	// t := T{}

	// err := yaml.Unmarshal([]byte(data), &t)
	// if err != nil {
	//         log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t:\n%v\n\n", t)

	a := AppConfig{}
	a.Users = append(a.Users, User{})
	a.Users = append(a.Users, User{})

	a.Users[0].Email = "user1@gmail.com"
	a.Users[1].Email = "asdf"

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

	err := validate.Struct(a)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		fmt.Println(errs)
	}

	d, err := yaml.Marshal(&a)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

}
