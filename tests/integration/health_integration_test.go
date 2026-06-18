package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/health/domain"
	"github.com/user/simple-blog/internal/platform/di"
	"github.com/user/simple-blog/internal/platform/server"
)

var _ = Describe("Health Integration", func() {
	var (
		cfg  *config.Config
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		srv  *server.Server
	)

	BeforeEach(func() {
		var err error
		cfg = &config.Config{
			App: config.AppConfig{
				Port: "8080",
			},
		}

		// Mock database
		dbRaw, mockRaw, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		Expect(err).NotTo(HaveOccurred())

		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw

		srv = di.InitializeServer(cfg, db)
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("GET /health", func() {
		Context("when database is reachable", func() {
			It("should return 200 OK with status UP", func() {
				// sqlmock expects a ping for health check
				mock.ExpectPing()

				req := httptest.NewRequest(http.MethodGet, "/health", nil)
				resp, err := srv.App.Test(req)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var fullResp struct {
					Success bool                  `json:"success"`
					Message string                `json:"message"`
					Data    domain.HealthResponse `json:"data"`
				}
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(body, &fullResp)
				Expect(err).NotTo(HaveOccurred())

				Expect(fullResp.Success).To(BeTrue())
				Expect(fullResp.Data.Status).To(Equal("UP"))
				Expect(fullResp.Data.Message).To(ContainSubstring("healthy"))

				Expect(mock.ExpectationsWereMet()).To(Succeed())
			})
		})

		Context("when database is NOT reachable", func() {
			It("should return 503 Service Unavailable with status DOWN", func() {
				mock.ExpectPing().WillReturnError(io.ErrUnexpectedEOF)

				req := httptest.NewRequest(http.MethodGet, "/health", nil)
				resp, err := srv.App.Test(req)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))

				var fullResp struct {
					Success bool                  `json:"success"`
					Message string                `json:"message"`
					Data    domain.HealthResponse `json:"data"`
				}
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(body, &fullResp)
				Expect(err).NotTo(HaveOccurred())

				Expect(fullResp.Success).To(BeFalse())
				Expect(fullResp.Data.Status).To(Equal("DOWN"))
				Expect(fullResp.Data.Message).To(ContainSubstring("failed"))

				Expect(mock.ExpectationsWereMet()).To(Succeed())
			})
		})
	})
})
