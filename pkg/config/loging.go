package config

type Logging struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	Color    bool   `envconfig:"LOG_COLOR" default:"true"`
	Text     bool   `envconfig:"LOG_TEXT" default:"true"`
	Pretty   bool   `envconfig:"LOG_PRETTY"`
}
