package user_repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/services/auth/entity"
)

type Repository interface {
	GetPosUser(ctx context.Context, id uuid.UUID) (*entity.PosUser, error)
	GetPosUserByEmail(ctx context.Context, email string) (*entity.PosUser, error)
	GetBackofficeUser(ctx context.Context, id uuid.UUID) (*entity.BackofficeUser, error)
	GetBackofficeUserByEmail(ctx context.Context, email string) (*entity.BackofficeUser, error)
	UpsertCustomer(ctx context.Context, data entity.Customer) error
}
