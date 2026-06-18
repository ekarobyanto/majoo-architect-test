package repository_test

import (
	"context"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/repository"
)

var _ = Describe("AuthRepository", func() {
	var (
		db     *sqlx.DB
		mock   sqlmock.Sqlmock
		repo   domain.AuthRepository
		ctx    context.Context
	)

	BeforeEach(func() {
		dbRaw, mockRaw, _ := sqlmock.New()
		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw
		repo = repository.NewAuthRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("GetByEmail", func() {
		It("should return a user when email exists", func() {
			email := "test@example.com"
			now := time.Now()
			rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
				AddRow("uuid-1", "testuser", email, "hash", now, now)

			mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
				WithArgs(email).
				WillReturnRows(rows)

			user, err := repo.GetByEmail(ctx, email)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			Expect(user.Email).To(Equal(email))
			Expect(user.Username).To(Equal("testuser"))
		})

		It("should return nil when user not found", func() {
			email := "notfound@example.com"
			mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
				WithArgs(email).
				WillReturnRows(sqlmock.NewRows([]string{"id"}))

			user, err := repo.GetByEmail(ctx, email)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).To(BeNil())
		})
	})

	Describe("GetUserRoles", func() {
		It("should return roles for a user", func() {
			userID := "user-uuid"
			now := time.Now()
			rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
				AddRow("role-1", "admin", now).
				AddRow("role-2", "user", now)

			mock.ExpectQuery("SELECT r.id, r.name, r.created_at FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = \\$1").
				WithArgs(userID).
				WillReturnRows(rows)

			roles, err := repo.GetUserRoles(ctx, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(roles).To(HaveLen(2))
			Expect(roles[0].Name).To(Equal("admin"))
			Expect(roles[1].Name).To(Equal("user"))
		})
	})
})
