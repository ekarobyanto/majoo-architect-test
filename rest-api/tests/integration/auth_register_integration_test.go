package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/platform/di"
	"github.com/user/simple-blog/internal/platform/server"
)

var _ = Describe("Auth Register Integration", func() {
	var (
		cfg  *config.Config
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		srv  *server.Server
	)

	BeforeEach(func() {
		cfg = &config.Config{
			App: config.AppConfig{
				Port: "8080",
			},
		}
		dbRaw, mockRaw, _ := sqlmock.New()
		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw
		srv = di.InitializeServer(cfg, db)
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("POST /auth/register", func() {
		It("should register a new user successfully", func() {
			reqBody := domain.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			}

			// Mock GetByUsername
			mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
				WithArgs(reqBody.Username).
				WillReturnRows(sqlmock.NewRows([]string{"id"}))

			// Mock GetByEmail
			mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
				WithArgs(reqBody.Email).
				WillReturnRows(sqlmock.NewRows([]string{"id"}))

			// Mock GetRoleByName
			mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = \\$1").
				WithArgs("user").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("role-uuid", "user"))

			// Transaction starts here
			mock.ExpectBegin()

			// Mock CreateUser
			now := time.Now()
			mock.ExpectQuery("INSERT INTO users").
				WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(now, now))

			// Mock AssignRole
			mock.ExpectExec("INSERT INTO user_roles").
				WithArgs(sqlmock.AnyArg(), "role-uuid").
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := srv.App.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			var fullResp struct {
				Success bool                    `json:"success"`
				Data    domain.RegisterResponse `json:"data"`
			}
			respBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(respBody, &fullResp)

			Expect(fullResp.Success).To(BeTrue())
			Expect(fullResp.Data.Username).To(Equal(reqBody.Username))
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
	})
})
