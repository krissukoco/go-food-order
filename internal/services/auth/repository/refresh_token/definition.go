package refresh_token_repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/auth"
	"github.com/krissukoco/go-food-order-microservices/internal/services/auth/entity"
)

type Repository interface {
	// Get gets refresh token entity from token.
	//
	// If exists, returns non-nil refresh token, else returns nil, entity.ErrNotFound.
	Get(ctx context.Context, token uuid.UUID) (*entity.RefreshToken, error)
	// Create creates a refresh token
	Create(ctx context.Context, expiredAt time.Time) (uuid.UUID, error)
	// DeleteAllOfUser deletes all refresh tokens of user
	DeleteAllOfUser(ctx context.Context, userId uuid.UUID, group auth.Group) error
}
