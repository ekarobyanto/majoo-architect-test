package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/go-backend-boilerplate/internal/modules/health/domain"
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
		return c.Status(fiber.StatusServiceUnavailable).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func RegisterRoutes(router fiber.Router, svc domain.HealthService) {
	h := NewHealthHandler(svc)
	router.Get("/health", h.CheckHealth)
}
