package config

import "github.com/spf13/viper"

// AppConfig holds general application configuration
type AppConfig struct {
	Port string `mapstructure:"PORT"`
}

// BindAppEnv binds application-related environment variables to Viper
func BindAppEnv(v *viper.Viper) {
	v.BindEnv("PORT")
}
