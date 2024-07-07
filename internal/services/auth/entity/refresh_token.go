package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/krissukoco/go-food-order-microservices/internal/auth"
)

type RefreshToken struct {
	Token     uuid.UUID  `json:"token"`
	UserId    uuid.UUID  `json:"user_id"`
	Group     auth.Group `json:"group"`
	ExpiredAt time.Time  `json:"expired_at"`
}
