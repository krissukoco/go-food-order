package customer_repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/services/user/entity"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
	GetByPhone(ctx context.Context, phone string) (*entity.Customer, error)
}
