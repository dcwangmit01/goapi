package config

import (
	"io/ioutil"
	"os"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "password"
	DefaultConfigFile    = "app.yaml"
)

var SingletonAppConfig *AppConfig

func init() {
	// Create the global AppConfig
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		// the file exists
		bytes, err := ioutil.ReadFile(DefaultConfigFile)
		if err != nil {
			panic("DefaultConfigFile exists but is not readable")
		}
		// parse
		SingletonAppConfig, _ = AppConfigFromYaml(string(bytes))
	} else {
		// Create a NewAppConfig
		SingletonAppConfig = NewAppConfig()
	}
}
