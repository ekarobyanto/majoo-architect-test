package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/user/simple-blog/config"
)

var openDB = sqlx.Open

// NewConnection initializes a new PostgreSQL connection using sqlx
func NewConnection(cfg *config.Config) (*sqlx.DB, error) {
	dsn := cfg.DB.URL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
		)
	} else if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		// lib/pq accepts URL DSNs directly; keep the value as-is.
	}

	db, err := openDB("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Connection pooling configuration
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
