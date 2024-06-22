package response

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	ProcessTime time.Duration `json:"process_time"`
	Data        interface{}   `json:"data"`
}

type ErrorResponse struct {
	HttpStatus int           `json:"-"`
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	Details    []ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func OK(c *fiber.Ctx, data interface{}, statusCode ...int) error {
	status := http.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	return c.Status(status).JSON(data)
}

func Error(c *fiber.Ctx, httpStatus int, code int, message string, details ...ErrorDetail) error {
	resp := ErrorResponse{
		HttpStatus: httpStatus,
		Code:       code,
		Message:    message,
		Details:    details,
	}
	if resp.Details == nil {
		resp.Details = make([]ErrorDetail, 0)
	}
	return c.Status(httpStatus).JSON(&resp)
}
