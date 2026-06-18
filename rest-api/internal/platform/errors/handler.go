package errors

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   *AppError `json:"error"`
}

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var appErr *AppError

	if e, ok := err.(*AppError); ok {
		code = e.HTTPStatus
		appErr = e
	} else if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		appErr = &AppError{
			Code:    CodeInternal,
			Message: e.Message,
		}
	} else {
		appErr = &AppError{
			Code:    CodeInternal,
			Message: err.Error(),
		}
	}

	return c.Status(code).JSON(ErrorResponse{
		Success: false,
		Error:   appErr,
	})
}
