package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func JSON(c *fiber.Ctx, status int, message string, data interface{}, meta ...interface{}) error {
	var m interface{}
	if len(meta) > 0 {
		m = meta[0]
	}
	return c.Status(status).JSON(Response{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
		Meta:    m,
	})
}

func Success(c *fiber.Ctx, status int, message string, data interface{}, meta ...interface{}) error {
	return JSON(c, status, message, data, meta...)
}
