package transport

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

func (x ErrorResponse) Error() string {
	return x.Message
}

func NewErrorResponse(code int, message string, details ...ErrorDetail) *ErrorResponse {
	e := &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
	if e.Details == nil {
		e.Details = make([]ErrorDetail, 0)
	}
	return e
}

type ErrorDetail struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func HTTPError(c *fiber.Ctx, httpStatus int, code int, message string, details ...ErrorDetail) error {
	e := NewErrorResponse(code, message, details...)
	return c.Status(httpStatus).JSON(e)
}

func GRPCError(grpcCode codes.Code, code int, message string, details ...ErrorDetail) error {
	e := NewErrorResponse(code, message, details...)
	b, _ := json.Marshal(e)
	return status.Error(grpcCode, string(b))
}

func ParseGRPCError(err error) (*ErrorResponse, codes.Code) {
	st, ok := status.FromError(err)
	if !ok {
		return &ErrorResponse{
			Code:    99998,
			Message: "Unknown: " + err.Error(),
		}, codes.Unknown
	}
	var e ErrorResponse
	if err := json.Unmarshal([]byte(st.Message()), &e); err != nil {
		return &ErrorResponse{
			Code:    99998,
			Message: "Unknown: " + err.Error(),
		}, st.Code()
	}
	return &e, st.Code()
}
