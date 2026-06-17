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
		defer os.Unsetenv("PORT")
		defer os.Unsetenv("DB_HOST")

		cfg, err := LoadConfig()
		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(cfg.Port).To(gomega.Equal("4000"))
		g.Expect(cfg.DB.Host).To(gomega.Equal("remote-host"))
	})
}
