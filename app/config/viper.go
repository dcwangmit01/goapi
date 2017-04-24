package config

import (
	"fmt"
)

var (
	ServerAddress string
	Host          = "localhost"
	Port          = 10080
)

func init() {
	ServerAddress = fmt.Sprintf("%s:%d", Host, Port)
}
