package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SendOtpResponse struct {
	Id        uuid.UUID `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}
