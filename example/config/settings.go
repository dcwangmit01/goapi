package config

type Settings struct {
	Initialized bool
	Debug       bool
	LogLevel    string `validate:"eq=debug|eq=info|eq=warning|eq=error|eq=fatal|eq=panic"`
}
