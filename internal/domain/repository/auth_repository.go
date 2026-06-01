//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../test/mock/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/google/uuid"
)

var (
	ErrAuthIdentityNotFound = errors.New("auth identity not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type AuthRepository interface {
	FindByProviderAndProviderUserID(
		ctx context.Context,
		provider entity.AuthProvider,
		providerUserID string,
	) (*entity.AuthIdentity, error)
	CreateAuthIdentity(ctx context.Context, authIdentity entity.AuthIdentity) error

	FindRefreshTokenByHash(
		ctx context.Context,
		tokenHash string,
	) (*entity.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error
	RevokeRefreshToken(
		ctx context.Context,
		refreshTokenID uuid.UUID,
		revokedAt time.Time,
	) error
	DeleteAuthIdentitiesByUserID(ctx context.Context, userID uuid.UUID) error
	UpdateRefreshTokenRevokedAtByUserID(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error
}
