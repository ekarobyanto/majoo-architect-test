package database

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/user/go-backend-boilerplate/config"
)

func TestNewConnection(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("should return error with invalid config", func(t *testing.T) {
		cfg := &config.Config{
			DB: config.DBConfig{
				Host: "invalid-host",
			},
		}

		db, err := NewConnection(cfg)
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(db).To(gomega.BeNil())
	})
}
