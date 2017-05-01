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

var AppConfig *appConfig

func init() {
	// Create the global AppConfig
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		// the file exists
		bytes, err := ioutil.ReadFile(DefaultConfigFile)
		if err != nil {
			panic("DefaultConfigFile exists but is not readable")
		}
		// parse
		AppConfig, _ = AppConfigFromYaml(string(bytes))
	} else {
		// Create a NewAppConfig
		AppConfig = NewAppConfig()
	}
}
