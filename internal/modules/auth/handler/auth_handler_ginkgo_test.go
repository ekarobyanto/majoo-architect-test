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
	"github.com/user/go-backend-boilerplate/internal/modules/auth/domain"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/handler"
	"github.com/user/go-backend-boilerplate/internal/platform/errors"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegisterResponse), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

var _ = Describe("AuthHandler", func() {
	var (
		app     *fiber.App
		mockSvc *mockAuthService
		h       *handler.AuthHandler
	)

	BeforeEach(func() {
		app = fiber.New(fiber.Config{
			ErrorHandler: errors.GlobalErrorHandler,
		})
		mockSvc = new(mockAuthService)
		h = handler.NewAuthHandler(mockSvc)
		app.Post("/auth/register", h.Register)
	})

	Describe("Register", func() {
		Context("with valid request", func() {
			It("should return 201 Created", func() {
				reqBody := domain.RegisterRequest{
					Username: "newuser",
					Email:    "new@example.com",
					Password: "password123",
				}
				expectedResp := &domain.RegisterResponse{
					ID:       "uuid-123",
					Username: "newuser",
					Email:    "new@example.com",
				}

				mockSvc.On("Register", mock.Anything, reqBody).Return(expectedResp, nil)

				jsonBody, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")
				
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var fullResp struct {
					Success bool                    `json:"success"`
					Data    domain.RegisterResponse `json:"data"`
				}
				body, _ := io.ReadAll(resp.Body)
				json.Unmarshal(body, &fullResp)
				
				Expect(fullResp.Success).To(BeTrue())
				Expect(fullResp.Data.Username).To(Equal("newuser"))
				mockSvc.AssertExpectations(GinkgoT())
			})
		})

		Context("with invalid request (missing fields)", func() {
			It("should return 422 Unprocessable Entity", func() {
				reqBody := `{"username": ""}`
				req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})
	})
})
