package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type RefreshToken struct {
	Token     uuid.UUID `json:"token"`
	UserId    int64     `json:"user_id"`
	ExpiredAt time.Time `json:"expired_at"`
}
