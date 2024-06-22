package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-food-order-microservices/internal/auth"
)

type contextKey struct {
	key string
}

var (
	authKey = &contextKey{"auth"}
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

func NewAuth(jwtSecret string) fiber.Handler {
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

		u, err := auth.ParseToken(split[1], jwtSecret)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrTokenExpired) {
				return fmt.Errorf("%w: %v", ErrUnauthorized, err)
			}
			return err
		}

		// set auth to locals
		c.Locals(authKey, u)

		return c.Next()
	}
}
