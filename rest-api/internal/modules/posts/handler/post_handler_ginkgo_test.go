package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/modules/posts/handler"
	"github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/models"
)

type mockPostService struct {
	mock.Mock
}

func (m *mockPostService) Create(ctx context.Context, authorID string, req domain.CreatePostRequest) (*models.Post, error) {
	args := m.Called(ctx, authorID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) GetByID(ctx context.Context, id string) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) GetDetailByID(ctx context.Context, id string) (*domain.PostDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PostDetailResponse), args.Error(1)
}

func (m *mockPostService) GetPaginated(ctx context.Context, query domain.PaginationQuery) (*domain.PaginatedPostResponse, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PaginatedPostResponse), args.Error(1)
}

func (m *mockPostService) Update(ctx context.Context, id string, user *authDomain.UserContext, req domain.UpdatePostRequest) (*models.Post, error) {
	args := m.Called(ctx, id, user, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) Delete(ctx context.Context, id string, user *authDomain.UserContext) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

var _ = Describe("PostHandler", func() {
	var (
		app     *fiber.App
		mockSvc *mockPostService
		h       *handler.PostHandler
	)

	BeforeEach(func() {
		app = fiber.New(fiber.Config{
			ErrorHandler: errors.GlobalErrorHandler,
		})
		mockSvc = new(mockPostService)
		h = handler.NewPostHandler(mockSvc)

		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user", &authDomain.UserContext{
				ID:    "user-1",
				Roles: []string{"user"},
			})
			return c.Next()
		})

		app.Post("/posts", h.Create)
		app.Get("/posts/:id", h.GetByID)
		app.Get("/posts", h.GetPaginated)
		app.Put("/posts/:id", h.Update)
		app.Delete("/posts/:id", h.Delete)
	})

	Describe("Create", func() {
		Context("with valid request", func() {
			It("should return 201 Created", func() {
				reqBody := domain.CreatePostRequest{
					Title:   "New Post",
					Content: "Content of new post",
				}
				expectedPost := &models.Post{
					ID:       "post-1",
					AuthorID: "user-1",
					Title:    reqBody.Title,
					Content:  reqBody.Content,
				}

				mockSvc.On("Create", mock.Anything, "user-1", reqBody).Return(expectedPost, nil)

				jsonBody, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var fullResp struct {
					Success bool        `json:"success"`
					Data    models.Post `json:"data"`
				}
				body, _ := io.ReadAll(resp.Body)
				json.Unmarshal(body, &fullResp)

				Expect(fullResp.Success).To(BeTrue())
				Expect(fullResp.Data.Title).To(Equal("New Post"))
				mockSvc.AssertExpectations(GinkgoT())
			})
		})

		Context("with invalid request (missing fields)", func() {
			It("should return 422 Unprocessable Entity", func() {
				reqBody := `{"title": ""}`
				req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})

		Context("when service returns an error", func() {
			It("should return mapped error status", func() {
				reqBody := domain.CreatePostRequest{Title: "New Post", Content: "Content"}
				mockSvc.On("Create", mock.Anything, "user-1", reqBody).Return(nil, errors.Internal("Failed to create post"))

				jsonBody, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GetByID", func() {
		It("should return 200 OK", func() {
			expected := &domain.PostDetailResponse{Post: models.Post{ID: "post-1", Title: "Title", Content: "Body"}}
			mockSvc.On("GetDetailByID", mock.Anything, "post-1").Return(expected, nil)

			req := httptest.NewRequest(http.MethodGet, "/posts/post-1", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return mapped error status when service fails", func() {
			mockSvc.On("GetDetailByID", mock.Anything, "post-1").Return(nil, errors.NotFound("Post not found"))

			req := httptest.NewRequest(http.MethodGet, "/posts/post-1", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GetPaginated", func() {
		It("should return 200 OK with paginated posts", func() {
			expected := &domain.PaginatedPostResponse{
				Data:       []models.Post{{ID: "post-1"}},
				Total:      1,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			}

			mockSvc.On("GetPaginated", mock.Anything, domain.PaginationQuery{Page: 1, Limit: 10}).Return(expected, nil)

			req := httptest.NewRequest(http.MethodGet, "/posts?page=1&limit=10", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return 400 for invalid query parameters", func() {
			req := httptest.NewRequest(http.MethodGet, "/posts?page=abc&limit=10", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should return mapped error status when service fails", func() {
			mockSvc.On("GetPaginated", mock.Anything, domain.PaginationQuery{Page: 1, Limit: 10}).
				Return(nil, errors.Internal("Failed to fetch posts"))

			req := httptest.NewRequest(http.MethodGet, "/posts?page=1&limit=10", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("Update", func() {
		It("should return 200 OK", func() {
			reqBody := domain.UpdatePostRequest{Title: "Updated title", Content: "Updated content"}
			mockSvc.On("Update", mock.Anything, "post-1", mock.AnythingOfType("*domain.UserContext"), reqBody).
				Return(&models.Post{ID: "post-1", Title: reqBody.Title, Content: reqBody.Content}, nil)

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/posts/post-1", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return 400 on invalid JSON body", func() {
			req := httptest.NewRequest(http.MethodPut, "/posts/post-1", strings.NewReader("{"))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should return mapped error status when service fails", func() {
			reqBody := domain.UpdatePostRequest{Title: "Updated title", Content: "Updated content"}
			mockSvc.On("Update", mock.Anything, "post-1", mock.AnythingOfType("*domain.UserContext"), reqBody).
				Return(nil, errors.Forbidden("You do not have permission to update this post"))

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPut, "/posts/post-1", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
		})
	})

	Describe("Delete", func() {
		It("should return 200 OK", func() {
			mockSvc.On("Delete", mock.Anything, "post-1", mock.AnythingOfType("*domain.UserContext")).Return(nil)

			req := httptest.NewRequest(http.MethodDelete, "/posts/post-1", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return mapped error status when service fails", func() {
			mockSvc.On("Delete", mock.Anything, "post-1", mock.AnythingOfType("*domain.UserContext")).
				Return(errors.Forbidden("You do not have permission to delete this post"))

			req := httptest.NewRequest(http.MethodDelete, "/posts/post-1", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
		})
	})
})
