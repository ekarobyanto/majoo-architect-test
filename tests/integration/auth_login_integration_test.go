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
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Auth Login Integration", func() {
	var (
		cfg  *config.Config
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		srv  *server.Server
	)

	BeforeEach(func() {
		cfg = &config.Config{
			Port: "8080",
			JWTSecret: "test-secret",
			JWTExpiration: 24,
		}
		dbRaw, mockRaw, _ := sqlmock.New()
		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw
		srv = di.InitializeServer(cfg, db)
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("POST /auth/login", func() {
		It("should login successfully and return token", func() {
			reqBody := domain.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
			
			// Mock GetByEmail
			mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
				WithArgs(reqBody.Email).
				WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
					AddRow("user-uuid", "testuser", reqBody.Email, string(hashedPassword), time.Now(), time.Now()))

			// Mock GetUserRoles
			mock.ExpectQuery("SELECT r.id, r.name, r.created_at FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = \\$1").
				WithArgs("user-uuid").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
					AddRow("role-uuid", "user", time.Now()))

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := srv.App.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var fullResp struct {
				Success bool                 `json:"success"`
				Data    domain.LoginResponse `json:"data"`
			}
			respBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(respBody, &fullResp)

			Expect(fullResp.Success).To(BeTrue())
			Expect(fullResp.Data.AccessToken).NotTo(BeEmpty())
			Expect(fullResp.Data.User.Email).To(Equal(reqBody.Email))
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
	})
})
