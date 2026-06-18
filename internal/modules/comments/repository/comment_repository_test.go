package repository_test

import (
	"context"
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/modules/comments/repository"
	"github.com/user/simple-blog/models"
)

var _ = Describe("CommentRepository", func() {
	var (
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		repo domain.CommentRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		dbRaw, mockRaw, _ := sqlmock.New()
		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw
		repo = repository.NewCommentRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("Create", func() {
		It("should successfully create a comment", func() {
			comment := &models.Comment{
				ID:       "comment-1",
				PostID:   "post-1",
				AuthorID: "author-1",
				Content:  "Test Content",
			}

			now := time.Now()
			rows := sqlmock.NewRows([]string{"created_at", "updated_at"}).
				AddRow(now, now)

			// sqlx NamedQuery rewrites :id etc. to $1, $2
			mock.ExpectQuery(`INSERT INTO comments \(id, post_id, author_id, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING created_at, updated_at`).
				WithArgs(comment.ID, comment.PostID, comment.AuthorID, comment.Content).
				WillReturnRows(rows)

			err := repo.Create(ctx, comment)
			Expect(err).NotTo(HaveOccurred())
			Expect(comment.CreatedAt).To(Equal(now))
			Expect(comment.UpdatedAt).To(Equal(now))
		})
	})

	Describe("GetByID", func() {
		It("should return a comment when found", func() {
			id := "comment-1"
			now := time.Now()
			rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "updated_at"}).
				AddRow(id, "post-1", "author-1", "Content", now, now)

			mock.ExpectQuery(`SELECT id, post_id, author_id, content, created_at, updated_at FROM comments WHERE id = \$1`).
				WithArgs(id).
				WillReturnRows(rows)

			comment, err := repo.GetByID(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(comment).NotTo(BeNil())
			Expect(comment.ID).To(Equal(id))
		})

		It("should return nil when comment not found", func() {
			id := "not-found"
			mock.ExpectQuery(`SELECT id, post_id, author_id, content, created_at, updated_at FROM comments WHERE id = \$1`).
				WithArgs(id).
				WillReturnError(sql.ErrNoRows)

			comment, err := repo.GetByID(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(comment).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should update a comment successfully", func() {
			comment := &models.Comment{
				ID:      "comment-1",
				Content: "Updated Content",
			}

			mock.ExpectExec(`UPDATE comments SET content = \$1, updated_at = \$2 WHERE id = \$3`).
				WithArgs(comment.Content, sqlmock.AnyArg(), comment.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.Update(ctx, comment)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Delete", func() {
		It("should delete a comment successfully", func() {
			id := "comment-1"
			mock.ExpectExec(`DELETE FROM comments WHERE id = \$1`).
				WithArgs(id).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.Delete(ctx, id)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
