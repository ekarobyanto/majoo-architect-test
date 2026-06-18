package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/modules/posts/repository"
	"github.com/user/simple-blog/models"
)

var _ = Describe("PostRepository", func() {
	var (
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		repo domain.PostRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		dbRaw, mockRaw, _ := sqlmock.New()
		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw
		repo = repository.NewPostRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("Create", func() {
		It("should successfully create a post", func() {
			post := &models.Post{
				ID:       "post-1",
				AuthorID: "author-1",
				Title:    "Test Title",
				Content:  "Test Content",
			}

			now := time.Now()
			rows := sqlmock.NewRows([]string{"created_at", "updated_at"}).
				AddRow(now, now)

			mock.ExpectQuery(`INSERT INTO posts \(id, author_id, title, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING created_at, updated_at`).
				WithArgs(post.ID, post.AuthorID, post.Title, post.Content).
				WillReturnRows(rows)

			err := repo.Create(ctx, post)
			Expect(err).NotTo(HaveOccurred())
			Expect(post.CreatedAt).To(Equal(now))
			Expect(post.UpdatedAt).To(Equal(now))
		})

		It("should return nil when insert returns no rows", func() {
			post := &models.Post{
				ID:       "post-2",
				AuthorID: "author-1",
				Title:    "No Row Title",
				Content:  "No Row Content",
			}

			rows := sqlmock.NewRows([]string{"created_at", "updated_at"})

			mock.ExpectQuery(`INSERT INTO posts \(id, author_id, title, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING created_at, updated_at`).
				WithArgs(post.ID, post.AuthorID, post.Title, post.Content).
				WillReturnRows(rows)

			err := repo.Create(ctx, post)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when insert fails", func() {
			post := &models.Post{
				ID:       "post-3",
				AuthorID: "author-1",
				Title:    "Error Title",
				Content:  "Error Content",
			}

			mock.ExpectQuery(`INSERT INTO posts \(id, author_id, title, content\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING created_at, updated_at`).
				WithArgs(post.ID, post.AuthorID, post.Title, post.Content).
				WillReturnError(errors.New("insert failed"))

			err := repo.Create(ctx, post)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insert failed"))
		})
	})

	Describe("GetByID", func() {
		It("should return a post when found", func() {
			id := "post-1"
			now := time.Now()
			rows := sqlmock.NewRows([]string{"id", "author_id", "title", "content", "created_at", "updated_at"}).
				AddRow(id, "author-1", "Title", "Content", now, now)

			mock.ExpectQuery(`SELECT id, author_id, title, content, created_at, updated_at FROM posts WHERE id = \$1`).
				WithArgs(id).
				WillReturnRows(rows)

			post, err := repo.GetByID(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(post).NotTo(BeNil())
			Expect(post.ID).To(Equal(id))
		})

		It("should return nil when post not found", func() {
			id := "not-found"
			mock.ExpectQuery(`SELECT id, author_id, title, content, created_at, updated_at FROM posts WHERE id = \$1`).
				WithArgs(id).
				WillReturnError(sql.ErrNoRows)

			post, err := repo.GetByID(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(post).To(BeNil())
		})
	})

	Describe("GetPaginated", func() {
		It("should return paginated posts and total count", func() {
			now := time.Now()
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
			mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).
				WillReturnRows(countRows)

			postRows := sqlmock.NewRows([]string{"id", "author_id", "title", "content", "created_at", "updated_at"}).
				AddRow("post-1", "author-1", "Title 1", "Content 1", now, now).
				AddRow("post-2", "author-2", "Title 2", "Content 2", now, now)

			mock.ExpectQuery(`SELECT id, author_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
				WithArgs(10, 0).
				WillReturnRows(postRows)

			posts, total, err := repo.GetPaginated(ctx, 1, 10)
			Expect(err).NotTo(HaveOccurred())
			Expect(total).To(Equal(int64(2)))
			Expect(posts).To(HaveLen(2))
		})

		It("should return error when count query fails", func() {
			mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).
				WillReturnError(errors.New("count failed"))

			posts, total, err := repo.GetPaginated(ctx, 1, 10)
			Expect(err).To(HaveOccurred())
			Expect(posts).To(BeNil())
			Expect(total).To(Equal(int64(0)))
			Expect(err.Error()).To(ContainSubstring("count failed"))
		})
	})

	Describe("Update", func() {
		It("should update a post successfully", func() {
			post := &models.Post{
				ID:      "post-1",
				Title:   "Updated Title",
				Content: "Updated Content",
			}

			mock.ExpectExec(`UPDATE posts SET title = \$1, content = \$2, updated_at = \$3 WHERE id = \$4`).
				WithArgs(post.Title, post.Content, sqlmock.AnyArg(), post.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.Update(ctx, post)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Delete", func() {
		It("should delete a post successfully", func() {
			id := "post-1"
			mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
				WithArgs(id).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.Delete(ctx, id)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
