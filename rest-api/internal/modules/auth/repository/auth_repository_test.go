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
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/repository"
	"github.com/user/simple-blog/models"
)

var _ = Describe("AuthRepository", func() {
	var (
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		repo domain.AuthRepository
		ctx  context.Context
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

	Describe("CreateUser", func() {
		It("should create user and populate timestamps", func() {
			now := time.Now()
			user := &models.User{
				ID:           "user-1",
				Username:     "newuser",
				Email:        "new@example.com",
				PasswordHash: "hash",
			}

			rows := sqlmock.NewRows([]string{"created_at", "updated_at"}).
				AddRow(now, now)

			mock.ExpectQuery("INSERT INTO users").
				WithArgs(user.ID, user.Username, user.Email, user.PasswordHash).
				WillReturnRows(rows)

			err := repo.CreateUser(ctx, user)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.CreatedAt).To(BeTemporally("~", now, time.Second))
			Expect(user.UpdatedAt).To(BeTemporally("~", now, time.Second))
		})

		It("should return query error", func() {
			user := &models.User{
				ID:           "user-2",
				Username:     "baduser",
				Email:        "bad@example.com",
				PasswordHash: "hash",
			}

			mock.ExpectQuery("INSERT INTO users").
				WithArgs(user.ID, user.Username, user.Email, user.PasswordHash).
				WillReturnError(errors.New("insert failed"))

			err := repo.CreateUser(ctx, user)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insert failed"))
		})
	})

	Describe("GetRoleByName", func() {
		It("should return role when found", func() {
			now := time.Now()
			rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
				AddRow("role-1", "user", now)

			mock.ExpectQuery("SELECT id, name, created_at FROM roles WHERE name = \\$1").
				WithArgs("user").
				WillReturnRows(rows)

			role, err := repo.GetRoleByName(ctx, "user")
			Expect(err).NotTo(HaveOccurred())
			Expect(role).NotTo(BeNil())
			Expect(role.ID).To(Equal("role-1"))
			Expect(role.Name).To(Equal("user"))
		})

		It("should return nil when role is not found", func() {
			mock.ExpectQuery("SELECT id, name, created_at FROM roles WHERE name = \\$1").
				WithArgs("missing").
				WillReturnError(sql.ErrNoRows)

			role, err := repo.GetRoleByName(ctx, "missing")
			Expect(err).NotTo(HaveOccurred())
			Expect(role).To(BeNil())
		})

		It("should return error when query fails", func() {
			mock.ExpectQuery("SELECT id, name, created_at FROM roles WHERE name = \\$1").
				WithArgs("user").
				WillReturnError(errors.New("query failed"))

			role, err := repo.GetRoleByName(ctx, "user")
			Expect(err).To(HaveOccurred())
			Expect(role).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("query failed"))
		})
	})

	Describe("AssignRole", func() {
		It("should assign role successfully", func() {
			mock.ExpectExec(`INSERT INTO user_roles \(user_id, role_id\) VALUES \(\$1, \$2\)`).
				WithArgs("user-1", "role-1").
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.AssignRole(ctx, "user-1", "role-1")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when insert fails", func() {
			mock.ExpectExec(`INSERT INTO user_roles \(user_id, role_id\) VALUES \(\$1, \$2\)`).
				WithArgs("user-1", "role-1").
				WillReturnError(errors.New("assign failed"))

			err := repo.AssignRole(ctx, "user-1", "role-1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("assign failed"))
		})
	})

	Describe("GetByUsername", func() {
		It("should return user when username exists", func() {
			now := time.Now()
			rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
				AddRow("uuid-1", "testuser", "test@example.com", "hash", now, now)

			mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
				WithArgs("testuser").
				WillReturnRows(rows)

			user, err := repo.GetByUsername(ctx, "testuser")
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			Expect(user.Username).To(Equal("testuser"))
		})

		It("should return nil when username is not found", func() {
			mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
				WithArgs("missing").
				WillReturnError(sql.ErrNoRows)

			user, err := repo.GetByUsername(ctx, "missing")
			Expect(err).NotTo(HaveOccurred())
			Expect(user).To(BeNil())
		})

		It("should return error when query fails", func() {
			mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
				WithArgs("testuser").
				WillReturnError(errors.New("db down"))

			user, err := repo.GetByUsername(ctx, "testuser")
			Expect(err).To(HaveOccurred())
			Expect(user).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("db down"))
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

		It("should return error when query fails", func() {
			userID := "user-uuid"
			mock.ExpectQuery("SELECT r.id, r.name, r.created_at FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = \\$1").
				WithArgs(userID).
				WillReturnError(errors.New("roles query failed"))

			roles, err := repo.GetUserRoles(ctx, userID)
			Expect(err).To(HaveOccurred())
			Expect(roles).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("roles query failed"))
		})
	})
})
