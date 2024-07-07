package jwt_usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krissukoco/go-food-order-microservices/internal/auth"
	"github.com/krissukoco/go-food-order-microservices/internal/services/auth/entity"
	refresh_token_repository "github.com/krissukoco/go-food-order-microservices/internal/services/auth/repository/refresh_token"
	user_repository "github.com/krissukoco/go-food-order-microservices/internal/services/auth/repository/user"
	"github.com/krissukoco/go-food-order-microservices/pkg/transaction"
)

type usecase struct {
	secret               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	transactioner        transaction.Transactioner
	refreshTokenRepo     refresh_token_repository.Repository
	userRepo             user_repository.Repository
}

var (
	jwtSigningMethod        = jwt.SigningMethodHS256
	issuer           string = "github.com/krissukoco/go-food-order"
	audience         string = "github.com/krissukoco/go-food-order"
)

func New(
	secret string,
	transactioner transaction.Transactioner,
	refreshTokenRepo refresh_token_repository.Repository,
	userRepo user_repository.Repository,
) auth.JWTAuthHandler {
	return &usecase{secret, 5 * time.Minute, 7 * 24 * time.Hour, transactioner, refreshTokenRepo, userRepo}
}

// Parse parses and validates access token and returns user
func (uc *usecase) Parse(ctx context.Context, accessToken string) (*auth.User, error) {
	var claims auth.JWTPayload
	t, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwtSigningMethod {
			return nil, fmt.Errorf("%w: signing method invalid", auth.ErrInvalidToken)
		}
		return []byte(uc.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}
	if !t.Valid {
		return nil, auth.ErrInvalidToken
	}
	if claims.Issuer != issuer {
		return nil, fmt.Errorf("%w: invalid issuer", auth.ErrInvalidToken)
	}
	var audienceValid bool
	for _, v := range claims.Audience {
		if v == audience {
			audienceValid = true
			break
		}
	}
	if !audienceValid {
		return nil, fmt.Errorf("%w: invalid audience", auth.ErrInvalidToken)
	}
	userId, err := uuid.FromString(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid subject", auth.ErrInvalidToken)
	}
	if !claims.Group.IsValid() {
		return nil, fmt.Errorf("%w: invalid user group", auth.ErrInvalidToken)
	}
	return &auth.User{
		Id:           userId,
		RestaurantId: claims.RestaurantId,
		Group:        claims.Group,
	}, nil
}

// GenerateTokens generate access and refresh tokens from user
func (uc *usecase) GenerateTokens(ctx context.Context, user auth.User) (access string, refresh string, err error) {
	now := time.Now()
	refreshTokenExp := now.Add(uc.refreshTokenDuration)
	refreshTokenUuid, err := uc.refreshTokenRepo.Create(ctx, refreshTokenExp)
	if err != nil {
		return
	}
	refresh = refreshTokenUuid.String()

	var jti string
	jti, err = uc.newJti()
	if err != nil {
		return
	}

	accessPayload := auth.JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   user.Id.String(),
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
		RestaurantId: user.RestaurantId,
		Group:        user.Group,
	}
	accessT := jwt.NewWithClaims(jwtSigningMethod, accessPayload)
	access, err = accessT.SignedString([]byte(uc.secret))
	if err != nil {
		return
	}

	return
}

// Refresh validates refresh token and returns new access and refresh tokens
func (uc *usecase) Refresh(ctx context.Context, refreshToken string) (access string, refresh string, err error) {
	// Get and validate refresh token
	refreshUuid, err := uuid.FromString(refreshToken)
	if err != nil {
		return "", "", auth.ErrInvalidToken
	}
	t, err := uc.refreshTokenRepo.Get(ctx, refreshUuid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return "", "", auth.ErrInvalidToken
		}
		return "", "", err
	}
	now := time.Now()
	if t.ExpiredAt.Before(now) {
		return "", "", auth.ErrTokenExpired
	}

	user := auth.User{
		Id:    t.UserId,
		Group: t.Group,
	}
	// Assert user still exists
	switch t.Group {
	case auth.Group_Backoffice:
		backOfficeUser, err := uc.userRepo.GetBackofficeUser(ctx, t.UserId)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				return "", "", auth.ErrTokenExpired
			}
			return "", "", nil
		}
		user.RestaurantId = backOfficeUser.RestaurantId
	case auth.Group_POS:
		posUser, err := uc.userRepo.GetPosUser(ctx, t.UserId)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				return "", "", auth.ErrTokenExpired
			}
			return "", "", nil
		}
		user.RestaurantId = posUser.RestaurantId
		// Note: customer does not have refresh token
	}

	// Delete and generate new tokens
	var accessToken string
	err = uc.transactioner.WithTx(ctx, func(c context.Context) error {
		err := uc.refreshTokenRepo.DeleteAllOfUser(c, user.Id, user.Group)
		if err != nil {
			return err
		}
		accessToken, refreshToken, err = uc.GenerateTokens(c, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (uc *usecase) newJti() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}
