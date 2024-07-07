package auth

import (
	"context"
	"errors"
)

type ctxKey struct {
	key string
}

var (
	authCtxKey = &ctxKey{"auth"}

	ErrUserContext = errors.New("user not found in context")
)

func FromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(authCtxKey).(*User)
	if !ok {
		return nil, ErrUserContext
	}
	return user, nil
}

func NewContext(parent context.Context, user *User) context.Context {
	return context.WithValue(parent, authCtxKey, user)
}
