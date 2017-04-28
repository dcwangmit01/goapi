package util

import (
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

var (
	validate = validator.New()
)

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

func StructToYamlStr(s interface{}) (string, error) {
	dump, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(dump), err
}

func StructFromYamlStr(s interface{}, yamlString string) error {
	err := yaml.Unmarshal([]byte(yamlString), s)
	return err
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
