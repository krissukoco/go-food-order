package login_usecase

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/services/auth/entity"
)

type Usecase interface {
	LoginPos(ctx context.Context, email, password string) (*entity.AuthTokens, error)
	LoginBackoffice(ctx context.Context, email, password string) (*entity.AuthTokens, error)
	SendCustomerOTP(ctx context.Context, phone string) (*entity.SendOtpResponse, error)
	VerifyCustomerOTP(ctx context.Context, otpId uuid.UUID, payload string) (*entity.AuthTokens, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)
