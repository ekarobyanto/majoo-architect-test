package config

import "github.com/spf13/viper"

// AuthConfig holds authentication-specific configuration
type AuthConfig struct {
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiration int    `mapstructure:"JWT_EXPIRATION_HOURS"`
}

// BindAuthEnv binds authentication-related environment variables to Viper
func BindAuthEnv(v *viper.Viper) {
	v.BindEnv("JWT_SECRET")
	v.BindEnv("JWT_EXPIRATION_HOURS")
	v.SetDefault("JWT_EXPIRATION_HOURS", 24)
}
