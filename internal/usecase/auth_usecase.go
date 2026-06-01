package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	domainauth "github.com/chocoko/character-support-walk-backend/internal/domain/auth"
	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/auth"
	"github.com/google/uuid"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)

type UserTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshOutput struct {
	AccessToken string `json:"access_token"`
}

type AuthUsecase struct {
	authRepo         repository.AuthRepository
	userRepo         repository.UserRepository
	characterRepo    repository.CharacterRepository
	jwtService       *auth.JWTService
	idTokenValidator domainauth.IDTokenValidator
	tx               repository.Transaction
}

func NewAuthUsecase(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	characterRepo repository.CharacterRepository,
	jwtService *auth.JWTService,
	idTokenValidator domainauth.IDTokenValidator,
	tx repository.Transaction,
) *AuthUsecase {
	return &AuthUsecase{
		authRepo:         authRepo,
		userRepo:         userRepo,
		characterRepo:    characterRepo,
		jwtService:       jwtService,
		idTokenValidator: idTokenValidator,
		tx:               tx,
	}
}

func (uc *AuthUsecase) LoginWithGoogle(ctx context.Context, idToken string) (UserTokens, error) {

	payload, err := uc.idTokenValidator.Validate(ctx, idToken)
	if err != nil {
		return UserTokens{}, fmt.Errorf("google id token subject is empty")
	}

	identity, err := uc.authRepo.FindByProviderAndProviderUserID(
		ctx,
		entity.AuthProviderGoogle,
		payload.Subject,
	)
	if err != nil {
		if !errors.Is(err, repository.ErrAuthIdentityNotFound) {
			return UserTokens{}, fmt.Errorf("failed to find auth identity: %w", err)
		}
		identity = nil
	}

	var (
		userID uuid.UUID
		output UserTokens
	)
	now := time.Now().UTC()

	if identity != nil {
		userID = identity.UserID
	} else {

		userID, err = uuid.NewV7()
		if err != nil {
			return UserTokens{}, fmt.Errorf("failed to generate user id: %w", err)
		}
		authIdentityID, err := uuid.NewV7()
		if err != nil {
			return UserTokens{}, fmt.Errorf("failed to generate auth identity id: %w", err)
		}

		err = uc.tx.RunInTx(ctx, func(ctx context.Context) error {

			user := entity.User{
				ID:        userID,
				Nickname:  "new user",
				CreatedAt: now,
				UpdatedAt: now,
			}
			err = uc.userRepo.Create(ctx, user)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			character := entity.Character{
				ID:            uuid.New(),
				UserID:        userID,
				SupporterType: entity.SupporterTypeDog,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			err = uc.characterRepo.Create(ctx, character)
			if err != nil {
				return fmt.Errorf("failed to create character: %w", err)
			}

			authIdentity := entity.AuthIdentity{
				ID:             authIdentityID,
				UserID:         userID,
				Provider:       entity.AuthProviderGoogle,
				ProviderUserID: payload.Subject,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			err = uc.authRepo.CreateAuthIdentity(ctx, authIdentity)
			if err != nil {
				return fmt.Errorf("failed to create auth identity: %w", err)
			}

			output, err = uc.issueTokens(ctx, userID, now)
			if err != nil {
				return fmt.Errorf("failed to issue tokens: %w", err)
			}
			return nil
		})

		if err != nil {
			return UserTokens{}, fmt.Errorf("failed to create new user login resources: %w", err)
		}
		return output, nil
	}

	return uc.issueTokens(ctx, userID, now)

}

func (uc *AuthUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (RefreshOutput, error) {
	tokenHash := hashRefreshToken(refreshToken)

	storedToken, err := uc.authRepo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("failed to find refresh token: %w", err)
	}
	if storedToken == nil {
		return RefreshOutput{}, ErrRefreshTokenNotFound
	}

	now := time.Now().UTC()

	if storedToken.RevokedAt != nil {
		return RefreshOutput{}, ErrRefreshTokenRevoked
	}

	if storedToken.ExpiresAt.Before(now) {
		return RefreshOutput{}, ErrRefreshTokenExpired
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(storedToken.UserID.String())
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	return RefreshOutput{
		AccessToken: accessToken,
	}, nil

}

func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {

	tokenHash := hashRefreshToken(refreshToken)

	storedToken, err := uc.authRepo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find refresh token: %w", err)
	}
	if storedToken == nil {
		return nil
	}

	now := time.Now().UTC()

	err = uc.authRepo.RevokeRefreshToken(ctx, storedToken.ID, now)
	if err != nil {
		return fmt.Errorf("failed to logout error: %w", err)
	}

	return nil

}

func (uc *AuthUsecase) issueTokens(ctx context.Context, userID uuid.UUID, now time.Time) (UserTokens, error) {
	accessToken, err := uc.jwtService.GenerateAccessToken(userID.String())
	if err != nil {
		return UserTokens{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return UserTokens{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTokenID, err := uuid.NewV7()
	if err != nil {
		return UserTokens{}, fmt.Errorf("failed to generate refresh token id: %w", err)
	}

	err = uc.authRepo.CreateRefreshToken(ctx, entity.RefreshToken{
		ID:        refreshTokenID,
		UserID:    userID,
		TokenHash: hashRefreshToken(refreshToken),
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		RevokedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return UserTokens{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
