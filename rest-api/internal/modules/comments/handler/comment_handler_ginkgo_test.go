package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/modules/comments/handler"
	"github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/models"
)

type mockCommentService struct {
	mock.Mock
}

func (m *mockCommentService) Create(ctx context.Context, postID, authorID string, req domain.CreateCommentRequest) (*models.Comment, error) {
	args := m.Called(ctx, postID, authorID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *mockCommentService) Update(ctx context.Context, id string, user *authDomain.UserContext, req domain.UpdateCommentRequest) (*models.Comment, error) {
	args := m.Called(ctx, id, user, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *mockCommentService) Delete(ctx context.Context, id string, user *authDomain.UserContext) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

var _ = Describe("CommentHandler", func() {
	var (
		app     *fiber.App
		mockSvc *mockCommentService
		h       *handler.CommentHandler
	)

	BeforeEach(func() {
		app = fiber.New(fiber.Config{
			ErrorHandler: errors.GlobalErrorHandler,
		})
		mockSvc = new(mockCommentService)
		h = handler.NewCommentHandler(mockSvc)

		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user", &authDomain.UserContext{
				ID:    "user-1",
				Roles: []string{"user"},
			})
			return c.Next()
		})

		app.Post("/posts/:id/comments", h.Create)
		app.Put("/comments/:id", h.Update)
		app.Delete("/comments/:id", h.Delete)
	})

	Describe("Create", func() {
		It("should return 201 Created", func() {
			reqBody := domain.CreateCommentRequest{Content: "New Comment"}
			mockSvc.On("Create", mock.Anything, "post-1", "user-1", reqBody).
				Return(&models.Comment{ID: "c1", Content: "New Comment"}, nil)

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/posts/post-1/comments", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		})

		It("should return 400 on invalid JSON body", func() {
			req := httptest.NewRequest(http.MethodPost, "/posts/post-1/comments", strings.NewReader("{"))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should return mapped error status when service fails", func() {
			reqBody := domain.CreateCommentRequest{Content: "New Comment"}
			mockSvc.On("Create", mock.Anything, "post-1", "user-1", reqBody).
				Return(nil, errors.Internal("Failed to create comment"))

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/posts/post-1/comments", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("Update", func() {
		It("should return 200 OK", func() {
			reqBody := domain.UpdateCommentRequest{Content: "Updated content"}
			mockSvc.On("Update", mock.Anything, "comment-1", mock.AnythingOfType("*domain.UserContext"), reqBody).
				Return(&models.Comment{ID: "comment-1", Content: reqBody.Content}, nil)

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/comments/comment-1", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return 400 on invalid JSON body", func() {
			req := httptest.NewRequest(http.MethodPut, "/comments/comment-1", strings.NewReader("{"))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should return mapped error status when service fails", func() {
			reqBody := domain.UpdateCommentRequest{Content: "Updated content"}
			mockSvc.On("Update", mock.Anything, "comment-1", mock.AnythingOfType("*domain.UserContext"), reqBody).
				Return(nil, errors.Forbidden("You do not have permission to update this comment"))

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/comments/comment-1", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
		})
	})

	Describe("Delete", func() {
		It("should return 200 OK", func() {
			mockSvc.On("Delete", mock.Anything, "comment-1", mock.AnythingOfType("*domain.UserContext")).Return(nil)

			req := httptest.NewRequest(http.MethodDelete, "/comments/comment-1", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return mapped error status when service fails", func() {
			mockSvc.On("Delete", mock.Anything, "comment-1", mock.AnythingOfType("*domain.UserContext")).
				Return(errors.Forbidden("You do not have permission to delete this comment"))

			req := httptest.NewRequest(http.MethodDelete, "/comments/comment-1", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
		})
	})
})
