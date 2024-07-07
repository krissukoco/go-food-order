package login_usecase

import (
	"context"
	"errors"

	"github.com/krissukoco/go-food-order-microservices/internal/services/auth/entity"
)

type Usecase interface {
	LoginPos(ctx context.Context, email, password string) (*entity.AuthTokens, error)
	LoginBackoffice(ctx context.Context, email, password string) (*entity.AuthTokens, error)
	LoginCustomer(ctx context.Context, phone string) error
	VerifyCustomerOTP(ctx context.Context, otp string) (*entity.AuthTokens, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)
