package common_usecase

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/services/user/entity"
)

type Usecase interface {
	GetPosUser(ctx context.Context, id uuid.UUID) (*entity.PosUser, error)
	GetPosUserByEmail(ctx context.Context, id uuid.UUID) (*entity.PosUser, error)
	GetBackofficeUser(ctx context.Context, id uuid.UUID) (*entity.BackofficeUser, error)
	GetBackofficeUserByEmail(ctx context.Context, email string) (*entity.BackofficeUser, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
}
