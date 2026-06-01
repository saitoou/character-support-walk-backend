package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/database"
	"github.com/google/uuid"
)

type AuthRepository struct {
	db database.DBTX
}

func NewAuthRepository(db database.DBTX) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByProviderAndProviderUserID(
	ctx context.Context,
	provider entity.AuthProvider,
	providerUserID string,
) (*entity.AuthIdentity, error) {

	query := `
	  SElECT
		id,
		user_id,
		provider,
		provider_user_id,
		created_at,
		updated_at
	  FROM auth_identities
	  WHERE provider_user_id = $1
	    AND provider = $2
	`

	var identity entity.AuthIdentity

	db := r.executor(ctx)
	row := db.QueryRowContext(ctx, query, providerUserID, provider)

	if err := row.Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderUserID,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrAuthIdentityNotFound
		}
		return nil, fmt.Errorf("failed to find identity: %w", err)
	}

	return &identity, nil
}

func (r *AuthRepository) CreateAuthIdentity(ctx context.Context, authIdentity entity.AuthIdentity) error {
	query := `
	  INSERT INTO auth_identities (
	    id,
		user_id,
		provider,
		provider_user_id,
		created_at,
		updated_at 
	  )
	  VALUES($1, $2, $3, $4, $5, $6)
	`

	db := r.executor(ctx)
	_, err := db.ExecContext(
		ctx,
		query,
		authIdentity.ID,
		authIdentity.UserID,
		authIdentity.Provider,
		authIdentity.ProviderUserID,
		authIdentity.CreatedAt,
		authIdentity.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert auth identity: %w", err)
	}

	return nil
}

func (r *AuthRepository) FindRefreshTokenByHash(
	ctx context.Context,
	tokenHash string,
) (*entity.RefreshToken, error) {

	query := `
	  SELECT
		id,
		user_id,
		token_hash,
		expires_at,
		revoked_at,
		created_at,
		updated_at 
	  FROM refresh_tokens
	  WHERE token_hash = $1
	`

	var refresh entity.RefreshToken

	db := r.executor(ctx)
	row := db.QueryRowContext(ctx, query, tokenHash)

	if err := row.Scan(
		&refresh.ID,
		&refresh.UserID,
		&refresh.TokenHash,
		&refresh.ExpiresAt,
		&refresh.RevokedAt,
		&refresh.CreatedAt,
		&refresh.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	return &refresh, nil
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error {
	query := `
	  INSERT INTO refresh_tokens (
	    id,
		user_id,
		token_hash,
		expires_at,
		revoked_at,
		created_at,
		updated_at 
	  )
	  VALUES($1, $2, $3, $4, $5, $6, $7)
	`

	db := r.executor(ctx)
	_, err := db.ExecContext(
		ctx,
		query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
		refreshToken.RevokedAt,
		refreshToken.CreatedAt,
		refreshToken.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert refresh token: %w", err)
	}

	return nil
}

func (r *AuthRepository) RevokeRefreshToken(
	ctx context.Context,
	refreshTokenID uuid.UUID,
	revokedAt time.Time,
) error {
	query := `
	  UPDATE refresh_tokens
	  SET
	    revoked_at = $1,
		updated_at = $2
	  WHERE id = $3
	`

	db := r.executor(ctx)
	_, err := db.ExecContext(ctx, query, revokedAt, revokedAt, refreshTokenID)
	if err != nil {
		return fmt.Errorf("failed to revork refresh token: %w", err)
	}

	return nil
}

func (r *AuthRepository) DeleteAuthIdentitiesByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
	  DELETE FROM auth_identities
	  WHERE user_id = $1;
	`
	db := r.executor(ctx)
	_, err := db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete auth identity: %w", err)
	}

	return nil
}

func (r *AuthRepository) UpdateRefreshTokenRevokedAtByUserID(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {

	query := `
	  UPDATE refresh_tokens
	  SET
	    revoked_at = $1,
		updated_at = $2
	  WHERE user_id = $3
	    AND revoked_at IS NULL;
	`
	db := r.executor(ctx)
	_, err := db.ExecContext(ctx, query, revokedAt, revokedAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update refresh token revoked_at: %w", err)
	}

	return nil
}

func (r *AuthRepository) executor(ctx context.Context) database.DBTX {
	return database.GetExecutor(ctx, r.db)
}
