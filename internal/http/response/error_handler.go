package response

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-food-order-microservices/internal/http/middleware"
	"github.com/krissukoco/go-food-order-microservices/internal/http/request"
	"github.com/sirupsen/logrus"
)

func NewErrorHandler(includeInternalError bool) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Unauthorized / Forbidden
		if errors.Is(err, middleware.ErrUnauthorized) {
			return Error(c, http.StatusUnauthorized, 10001, err.Error())
		}
		if errors.Is(err, middleware.ErrForbidden) {
			return Error(c, http.StatusForbidden, 10002, err.Error())
		}

		// Request body error
		if errors.Is(err, request.ErrInvalidRequestBody) {
			return c.Status(http.StatusUnprocessableEntity).JSON(&ErrorResponse{
				Code:    29999,
				Message: err.Error(),
				Details: []ErrorDetail{
					{Field: "body", Message: "Unprocessible request body"},
				},
			})
		}

		// Validation error
		var valErr *request.ValidationError
		if errors.As(err, &valErr) {
			details := make([]ErrorDetail, len(valErr.Errors))
			for i, v := range valErr.Errors {
				details[i] = ErrorDetail{Field: v.Field, Message: v.Message}
			}
			return c.Status(http.StatusUnprocessableEntity).JSON(&ErrorResponse{
				Code:    29998,
				Message: valErr.Message,
				Details: details,
			})
		}

		// Fiber error
		var fibErr *fiber.Error
		if errors.As(err, &fibErr) {
			return Error(c, fibErr.Code, fibErr.Code, fibErr.Message)
		}

		// Unknown / Internal Error
		logrus.WithError(err).Error("UNKNOWN INTERNAL ERROR")
		details := make([]ErrorDetail, 0)
		if includeInternalError {
			details = append(details, ErrorDetail{
				Field:   "INTERNAL_ERROR_MESSAGE",
				Message: err.Error(),
			})
		}
		return c.Status(500).JSON(&ErrorResponse{
			Code:    99999,
			Message: "Internal server error",
			Details: details,
		})
	}
}
