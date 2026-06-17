package config

import "github.com/spf13/viper"

// DBConfig holds the database-specific configuration
type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

// BindDBEnv binds database-related environment variables to Viper
func BindDBEnv(v *viper.Viper) {
	v.BindEnv("DB_HOST")
	v.BindEnv("DB_PORT")
	v.BindEnv("DB_USER")
	v.BindEnv("DB_PASSWORD")
	v.BindEnv("DB_NAME")
	v.BindEnv("DB_SSLMODE")
}
