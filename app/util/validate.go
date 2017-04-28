package util

import (
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

func init() {
	// Register an additional validator called "phone"
	// A phone number contains no dashes, spaces, and parens. Valid keys
	// all keypad [0-9#*] as well as "+" to designate country code, and "w"
	// to signify a .5 second wait.
	var phoneRegex = regexp.MustCompile(`^\+?[0-9#*w]+$`)
	var isPhone = func(fl validator.FieldLevel) bool {
		return phoneRegex.MatchString(fl.Field().String())
	}

	validate = validator.New()
	validate.RegisterValidation("phone", isPhone)
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func ValidationErrorToMap(err error) map[string]string {
	// For test validation only.  Returns empty map if no errors.  Else
	// returns a map with key=namespace, value=validation_tag_that_failed

	m := make(map[string]string)
	if err == nil {
		return m
	}

	// translate the validationErrors above from odd types to
	// map[string]string
	validationErrors := err.(validator.ValidationErrors)
	for i := 0; i < len(validationErrors); i++ {
		fe := validationErrors[i]
		m[fe.Namespace()] = fe.Tag()
	}
	return m
}
