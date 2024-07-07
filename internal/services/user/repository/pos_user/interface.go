package pos_user_repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/services/user/entity"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.PosUser, error)
	GetByEmail(ctx context.Context, email string) (*entity.PosUser, error)
}
