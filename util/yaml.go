package util

import (
	yaml "gopkg.in/yaml.v2"
)

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
