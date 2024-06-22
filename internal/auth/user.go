package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has been expired")
	ErrInvalidUserGroup = errors.New("invalid user group")
)

type Group string

var _ json.Marshaler = (*Group)(nil)   // transform to JSON
var _ json.Unmarshaler = (*Group)(nil) // helps with validating user group

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

func (x Group) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, x.String())), nil
}

func (x *Group) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}
	g := Group(s)
	if !g.IsValid() {
		return fmt.Errorf("invalid user group '%s'", s)
	}
	*x = g
	return nil
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

func (j JWTPayload) User() (*User, error) {
	userId, err := uuid.FromString(j.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return &User{Id: uuid.UUID(userId), Group: j.Group, RestaurantId: j.RestaurantId}, nil
}

func ParseToken(token string, secret string) (*User, error) {
	var p JWTPayload
	t, err := jwt.ParseWithClaims(token, &p, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, err
	}
	return p.User()
}
