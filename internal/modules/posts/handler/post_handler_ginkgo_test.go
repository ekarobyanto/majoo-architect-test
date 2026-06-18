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
	})
})
