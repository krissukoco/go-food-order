package transport

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-food-order-microservices/internal/auth"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

func NewAuthMiddleware(hdl auth.JWTAuthHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("authorization")
		if h == "" {
			return ErrUnauthorized
		}
		split := strings.Split(h, " ")
		if len(split) != 2 {
			return ErrUnauthorized
		}

		// only allow 'bearer' method
		if strings.ToLower(split[0]) != "bearer" {
			return ErrUnauthorized
		}

		u, err := hdl.Parse(c.Context(), split[1])
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrTokenExpired) {
				return fmt.Errorf("%w: %v", ErrUnauthorized, err)
			}
			return err
		}

		// set auth user to context
		ctx := auth.NewContext(c.Context(), u)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func HTTPAuthContext(c *fiber.Ctx) (*auth.User, error) {
	return auth.FromContext(c.UserContext())
}
