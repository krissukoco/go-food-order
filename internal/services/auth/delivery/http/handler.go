package auth_http_handler

import (
	"github.com/gofiber/fiber/v2"
	auth_http_api "github.com/krissukoco/go-food-order-microservices/api/openapi/auth"
	login_usecase "github.com/krissukoco/go-food-order-microservices/internal/services/auth/usecase/auth/login"
)

type handler struct {
	logicUc login_usecase.Usecase
}

var _ auth_http_api.ServerInterface = (*handler)(nil)

func (h *handler) LoginPos(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *handler) LoginBackoffice(c *fiber.Ctx) error {
	panic("unimplemented")
}
