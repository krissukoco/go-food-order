package transport

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

var (
	val = newValidator()
)

type HTTPRequest struct {
	Params interface{}
	Query  interface{}
	Body   interface{}
}

func ParseAndValidateHTTPRequest(c *fiber.Ctx, req *HTTPRequest) error {
	if req.Params != nil {
		if err := c.ParamsParser(req.Params); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidRequest, err)
		}
		if err := val.Struct(req.Params); err != nil {
			return handleRequestValidationError(err)
		}
	}
	if req.Query != nil {
		if err := c.QueryParser(req.Query); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidRequest, err)
		}
		if err := val.Struct(req.Query); err != nil {
			return handleRequestValidationError(err)
		}
	}
	if req.Body != nil {
		if err := c.BodyParser(req.Body); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidRequest, err)
		}
		if err := val.Struct(req.Body); err != nil {
			return handleRequestValidationError(err)
		}
	}
	return nil
}

type ValidationError struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

func (v ValidationError) Error() string {
	return v.Message
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func newValidator() *validator.Validate {
	val := validator.New(validator.WithRequiredStructEnabled())
	val.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			return field.Name
		}
		split := strings.Split(jsonTag, ",")
		jsonFieldName := split[0]
		if jsonFieldName == "-" {
			return field.Name
		}
		return jsonFieldName
	})
	return val
}

func handleRequestValidationError(err error) error {
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		errMsg := "Validation Error"
		fieldErrors := make([]FieldError, len(valErrs))
		for i, v := range valErrs {
			fieldErrors[i] = FieldError{
				Field:   v.Tag(),
				Message: fmt.Sprintf("Failed on tag '%s': %s", v.Tag(), v.ActualTag()),
			}
		}
		if len(fieldErrors) > 0 {
			errMsg += ": " + fieldErrors[0].Message
		}
		return &ValidationError{
			Message: errMsg,
			Errors:  fieldErrors,
		}
	}
	return err
}
