package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App  AppConfig  `mapstructure:",squash"`
	Auth AuthConfig `mapstructure:",squash"`
	DB   DBConfig   `mapstructure:",squash"`
}

// LoadConfig loads configuration from .env file or environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind env vars via modular functions
	BindAppEnv(v)
	BindAuthEnv(v)
	BindDBEnv(v)

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	// Ignore error if config file not found, fallback to env vars
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
