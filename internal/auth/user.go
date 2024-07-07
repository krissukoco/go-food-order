package auth

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has been expired")
	ErrInvalidUserGroup = errors.New("invalid user group")
)

type Group string

func (x Group) String() string {
	return string(x)
}

func (x Group) IsValid() bool {
	switch x {
	case Group_Customer, Group_POS, Group_Backoffice:
		return true
	}
	return false
}

const (
	Group_Customer   Group = "CUSTOMER"
	Group_POS        Group = "POS"
	Group_Backoffice Group = "BACKOFFICE"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Group        Group     `json:"group"`
	RestaurantId uuid.UUID `json:"restaurant_id"`
}

type JWTPayload struct {
	jwt.RegisteredClaims
	Group        Group     `json:"ug"`
	RestaurantId uuid.UUID `json:"rid,omitempty"`
}

type JWTAuthHandler interface {
	// Parse parses and validates access token and returns user
	Parse(ctx context.Context, accessToken string) (*User, error)
	// GenerateTokens generate access and refresh tokens from user
	GenerateTokens(ctx context.Context, user User) (access string, refresh string, err error)
	// Refresh validates refresh token and returns new access and refresh tokens
	Refresh(ctx context.Context, refreshToken string) (access string, refresh string, err error)
}
