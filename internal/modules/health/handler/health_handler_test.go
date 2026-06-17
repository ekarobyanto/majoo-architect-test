package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-backend-boilerplate/internal/modules/health/domain"
	"github.com/user/go-backend-boilerplate/internal/modules/health/handler"
)

type mockHealthService struct {
	mock.Mock
}

func (m *mockHealthService) Check(ctx context.Context) (domain.HealthResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.HealthResponse), args.Error(1)
}

func TestHealthHandler_CheckHealth(t *testing.T) {
	app := fiber.New()
	mockSvc := new(mockHealthService)
	h := handler.NewHealthHandler(mockSvc)
	
	app.Get("/health", h.CheckHealth)

	expectedResponse := domain.HealthResponse{
		Status:  "UP",
		Message: "Database connection is healthy",
	}

	mockSvc.On("Check", mock.Anything).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fullResp struct {
		Success bool                  `json:"success"`
		Message string                `json:"message"`
		Data    domain.HealthResponse `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &fullResp)

	assert.True(t, fullResp.Success)
	assert.Equal(t, expectedResponse, fullResp.Data)
	mockSvc.AssertExpectations(t)
}
