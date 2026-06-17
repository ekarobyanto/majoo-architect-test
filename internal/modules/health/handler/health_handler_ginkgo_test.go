package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-backend-boilerplate/internal/modules/health/domain"
	"github.com/user/go-backend-boilerplate/internal/modules/health/handler"
)

var _ = Describe("HealthHandler", func() {
	var (
		app     *fiber.App
		mockSvc *mockHealthService
		h       *handler.HealthHandler
	)

	BeforeEach(func() {
		app = fiber.New()
		mockSvc = new(mockHealthService)
		h = handler.NewHealthHandler(mockSvc)
		app.Get("/health", h.CheckHealth)
	})

	Describe("CheckHealth", func() {
		Context("when the health service returns UP", func() {
			It("should return 200 OK and the health response", func() {
				expectedResponse := domain.HealthResponse{
					Status:  "UP",
					Message: "Database connection is healthy",
				}

				mockSvc.On("Check", mock.Anything).Return(expectedResponse, nil)

				req := httptest.NewRequest(http.MethodGet, "/health", nil)
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var fullResp struct {
					Success bool                  `json:"success"`
					Message string                `json:"message"`
					Data    domain.HealthResponse `json:"data"`
				}
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(body, &fullResp)
				Expect(err).NotTo(HaveOccurred())

				Expect(fullResp.Success).To(BeTrue())
				Expect(fullResp.Data).To(Equal(expectedResponse))
				mockSvc.AssertExpectations(GinkgoT())
			})
		})

		Context("when the health service returns DOWN", func() {
			It("should return 503 Service Unavailable and the health response", func() {
				expectedResponse := domain.HealthResponse{
					Status:  "DOWN",
					Message: "Database connection is unhealthy",
				}

				mockSvc.On("Check", mock.Anything).Return(expectedResponse, context.DeadlineExceeded)

				req := httptest.NewRequest(http.MethodGet, "/health", nil)
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))

				var fullResp struct {
					Success bool                  `json:"success"`
					Message string                `json:"message"`
					Data    domain.HealthResponse `json:"data"`
				}
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(body, &fullResp)
				Expect(err).NotTo(HaveOccurred())

				Expect(fullResp.Success).To(BeFalse())
				Expect(fullResp.Data).To(Equal(expectedResponse))
				mockSvc.AssertExpectations(GinkgoT())
			})
		})
	})
})
