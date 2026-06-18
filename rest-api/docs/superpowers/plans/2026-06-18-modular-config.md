# Modularize Configuration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the configuration system to separate App, Auth, and Database configurations into their own files and nested structures for better modularity.

**Architecture:** We will create `AppConfig` and `AuthConfig` structs in separate files, following the pattern established by `DBConfig`. The main `Config` struct will embed these sub-structs. Binding and default logic will also be moved to the respective modular files.

**Tech Stack:** Go, Viper.

---

### Task 1: Create App and Auth Configuration Files

**Files:**
- Create: `config/app_config.go`
- Create: `config/auth_config.go`

- [ ] **Step 1: Create `config/app_config.go`**

```go
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
```

- [ ] **Step 2: Create `config/auth_config.go`**

```go
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
```

- [ ] **Step 3: Commit**

```bash
git add config/app_config.go config/auth_config.go
git commit -m "feat(config): add modular app and auth configuration files"
```

---

### Task 2: Refactor Main Configuration

**Files:**
- Modify: `config/config.go`

- [ ] **Step 1: Update `Config` struct and `LoadConfig`**

```go
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
```

- [ ] **Step 2: Commit**

```bash
git add config/config.go
git commit -m "refactor(config): use modular config structs in main config"
```

---

### Task 3: Update Usage in Codebase

**Files:**
- Modify: `internal/modules/auth/service/auth_service.go`
- Modify: `internal/modules/auth/service/auth_service_test.go`
- Modify: `config/config_test.go`
- Modify: `tests/integration/auth_login_integration_test.go`
- Modify: `tests/integration/auth_register_integration_test.go`

- [ ] **Step 1: Update `internal/modules/auth/service/auth_service.go`**

Change `s.cfg.JWTExpiration` to `s.cfg.Auth.JWTExpiration` and `s.cfg.JWTSecret` to `s.cfg.Auth.JWTSecret`.

- [ ] **Step 2: Update `internal/modules/auth/service/auth_service_test.go`**

Update mock config initialization.

- [ ] **Step 3: Update `config/config_test.go`**

Update assertions to reflect the new structure.

- [ ] **Step 4: Update Integration Tests**

Update `tests/integration/auth_login_integration_test.go` and `tests/integration/auth_register_integration_test.go`.

- [ ] **Step 5: Run all tests**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "refactor: update config usage to use modular structure"
```
