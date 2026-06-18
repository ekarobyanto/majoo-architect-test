package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Port          string   `mapstructure:"PORT"`
	JWTSecret     string   `mapstructure:"JWT_SECRET"`
	JWTExpiration int      `mapstructure:"JWT_EXPIRATION_HOURS"`
	DB            DBConfig `mapstructure:",squash"`
}

// LoadConfig loads configuration from .env file or environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind general env vars
	v.BindEnv("PORT")
	v.BindEnv("JWT_SECRET")
	v.BindEnv("JWT_EXPIRATION_HOURS")

	// Set defaults
	v.SetDefault("JWT_EXPIRATION_HOURS", 24)

	// Bind module-specific env vars
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
