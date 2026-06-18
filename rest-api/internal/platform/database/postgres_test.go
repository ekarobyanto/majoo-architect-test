package database

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/user/simple-blog/config"
)

// anyTime is a helper for matching time.Duration
type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	return true
}

func TestNewConnection(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("should return error with invalid config", func(t *testing.T) {
		cfg := &config.Config{
			DB: config.DBConfig{
				Host: "invalid-host",
			},
		}

		// Ensure we use the real postgres driver which will fail Ping
		oldOpenDB := openDB
		defer func() { openDB = oldOpenDB }()
		// No need to change openDB here as it defaults to sqlx.Open

		db, err := NewConnection(cfg)
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(db).To(gomega.BeNil())
	})

	t.Run("should return success with valid config using mock", func(t *testing.T) {
		dbMock, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		g.Expect(err).NotTo(gomega.HaveOccurred())
		defer dbMock.Close()

		mock.ExpectPing()

		// Temporarily swap the openDB function
		oldOpenDB := openDB
		openDB = func(driver, dsn string) (*sqlx.DB, error) {
			return sqlx.NewDb(dbMock, "sqlmock"), nil
		}
		defer func() { openDB = oldOpenDB }()

		cfg := &config.Config{
			DB: config.DBConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "user",
				Password: "password",
				Name:     "db",
				SSLMode:  "disable",
			},
		}

		db, err := NewConnection(cfg)
		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(db).NotTo(gomega.BeNil())
		g.Expect(mock.ExpectationsWereMet()).To(gomega.Succeed())
	})

	t.Run("should return error if openDB fails", func(t *testing.T) {
		// Temporarily swap the openDB function to return error
		oldOpenDB := openDB
		openDB = func(drv, dsn string) (*sqlx.DB, error) {
			return nil, driver.ErrBadConn
		}
		defer func() { openDB = oldOpenDB }()

		cfg := &config.Config{
			DB: config.DBConfig{
				Host: "localhost",
			},
		}

		db, err := NewConnection(cfg)
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(db).To(gomega.BeNil())
	})
}
