package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/go-backend-boilerplate/internal/modules/health/domain"
	"github.com/user/go-backend-boilerplate/internal/platform/response"
	"net/http"
)

type HealthHandler struct {
	svc domain.HealthService
}

func NewHealthHandler(svc domain.HealthService) *HealthHandler {
	return &HealthHandler{
		svc: svc,
	}
}

// CheckHealth godoc
// @Summary Check service health
// @Description Get the status of the service and database
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} domain.HealthResponse
// @Router /health [get]
func (h *HealthHandler) CheckHealth(c *fiber.Ctx) error {
	resp, err := h.svc.Check(c.Context())
	if err != nil {
		return response.JSON(c, http.StatusServiceUnavailable, "Service Unavailable", resp)
	}

	return response.Success(c, http.StatusOK, "Health check success", resp)
}

func RegisterRoutes(router fiber.Router, h *HealthHandler) {
	router.Get("/health", h.CheckHealth)
}
