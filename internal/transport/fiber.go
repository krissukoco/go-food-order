package transport

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewFiberErrorHandler(includeInternalError bool) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Unauthorized / Forbidden
		if errors.Is(err, ErrUnauthorized) {
			return HTTPError(c, http.StatusUnauthorized, 10001, err.Error())
		}
		if errors.Is(err, ErrForbidden) {
			return HTTPError(c, http.StatusForbidden, 10002, err.Error())
		}

		// Request body error
		if errors.Is(err, ErrInvalidRequest) {
			return HTTPError(c, http.StatusUnprocessableEntity, 29999, err.Error())
		}

		// Validation error
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			details := make([]ErrorDetail, len(valErr.Errors))
			for i, v := range valErr.Errors {
				details[i] = ErrorDetail{Field: v.Field, Error: v.Message}
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
			return HTTPError(c, fibErr.Code, fibErr.Code, fibErr.Message)
		}

		// Unknown / Internal Error
		logrus.WithError(err).Error("UNKNOWN INTERNAL ERROR")
		details := make([]ErrorDetail, 0)
		if includeInternalError {
			details = append(details, ErrorDetail{
				Field: "INTERNAL_ERROR_MESSAGE",
				Error: err.Error(),
			})
		}
		return c.Status(500).JSON(&ErrorResponse{
			Code:    99999,
			Message: "Internal server error",
			Details: details,
		})
	}
}
