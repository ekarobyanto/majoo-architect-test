package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/platform/di"
	"github.com/user/simple-blog/internal/platform/server"
	"github.com/user/simple-blog/models"
)

func generateTestToken(secret, sub string, roles []string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"roles": roles,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

var _ = Describe("Posts Integration", func() {
	var (
		cfg  *config.Config
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		srv  *server.Server
	)

	BeforeEach(func() {
		cfg = &config.Config{
			App: config.AppConfig{Port: "8080"},
			Auth: config.AuthConfig{
				JWTSecret: "test-secret",
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

	Describe("POST /posts", func() {
		It("should create a post when authenticated", func() {
			reqBody := domain.CreatePostRequest{
				Title:   "Integration Post",
				Content: "Content of integration post",
			}

			mock.ExpectQuery("INSERT INTO posts").
				WithArgs(sqlmock.AnyArg(), "user-1", reqBody.Title, reqBody.Content).
				WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
					AddRow(time.Now(), time.Now()))

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+generateTestToken(cfg.Auth.JWTSecret, "user-1", []string{"user"}))

			resp, err := srv.App.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			var fullResp struct {
				Success bool        `json:"success"`
				Data    models.Post `json:"data"`
			}
			respBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(respBody, &fullResp)

			Expect(fullResp.Success).To(BeTrue())
			Expect(fullResp.Data.Title).To(Equal("Integration Post"))
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
	})
})
