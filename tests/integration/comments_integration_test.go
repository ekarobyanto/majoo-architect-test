package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/platform/di"
	"github.com/user/simple-blog/internal/platform/server"
	"github.com/user/simple-blog/models"
)

var _ = Describe("Comments Integration", func() {
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

	Describe("POST /posts/:id/comments", func() {
		It("should add a comment when authenticated", func() {
			postID := "post-1"
			userID := "user-1"
			reqBody := domain.CreateCommentRequest{Content: "Great post!"}

			// Mock post existence check
			mock.ExpectQuery("SELECT (.+) FROM posts WHERE id = \\$1").
				WithArgs(postID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(postID, "author-1"))

			// Mock comment creation
			mock.ExpectQuery("INSERT INTO comments").
				WithArgs(sqlmock.AnyArg(), postID, userID, reqBody.Content).
				WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
					AddRow(time.Now(), time.Now()))

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/posts/"+postID+"/comments", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+generateTestToken(cfg.Auth.JWTSecret, userID, []string{"user"}))

			resp, err := srv.App.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			var fullResp struct {
				Success bool           `json:"success"`
				Data    models.Comment `json:"data"`
			}
			json.NewDecoder(resp.Body).Decode(&fullResp)
			Expect(fullResp.Success).To(BeTrue())
			Expect(fullResp.Data.Content).To(Equal("Great post!"))
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
	})
})
