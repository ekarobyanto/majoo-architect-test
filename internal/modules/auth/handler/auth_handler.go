package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/platform/response"
	"github.com/user/simple-blog/internal/platform/validation"
)

type AuthHandler struct {
	svc domain.AuthService
}

func NewAuthHandler(svc domain.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and assign default role
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=domain.RegisterResponse}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	resp, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, "Registration successful", resp)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=domain.LoginResponse}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	resp, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Login successful", resp)
}
