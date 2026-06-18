package config

import (
	"os"
	"testing"

	"github.com/onsi/gomega"
)

func TestLoadConfig(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("should load from environment variables", func(t *testing.T) {
		os.Setenv("PORT", "4000")
		os.Setenv("DB_HOST", "remote-host")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("JWT_EXPIRATION_HOURS", "48")
		defer os.Unsetenv("PORT")
		defer os.Unsetenv("DB_HOST")
		defer os.Unsetenv("JWT_SECRET")
		defer os.Unsetenv("JWT_EXPIRATION_HOURS")

		cfg, err := LoadConfig()
		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(cfg.App.Port).To(gomega.Equal("4000"))
		g.Expect(cfg.DB.Host).To(gomega.Equal("remote-host"))
		g.Expect(cfg.Auth.JWTSecret).To(gomega.Equal("test-secret"))
		g.Expect(cfg.Auth.JWTExpiration).To(gomega.Equal(48))
	})

	t.Run("should use default JWT expiration", func(t *testing.T) {
		cfg, err := LoadConfig()
		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(cfg.Auth.JWTExpiration).To(gomega.Equal(24))
	})
}
