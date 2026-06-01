package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domainauth "github.com/chocoko/character-support-walk-backend/internal/domain/auth"
	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/auth"
	authmock "github.com/chocoko/character-support-walk-backend/internal/test/mock/auth"
	mock "github.com/chocoko/character-support-walk-backend/internal/test/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AuthUsecase_RefreshAccessToken(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	refreshToken := "test-refresh-token"
	tokenHash := hashRefreshToken(refreshToken)
	now := time.Now().UTC().Truncate(time.Millisecond)
	findErr := errors.New("failed to find refresh token")
	jwtService := auth.NewJWTService("test-secret", time.Hour, "test-issuer")

	tests := []struct {
		name                  string
		mockFunc              func(*mock.MockAuthRepository)
		wantAccessTokenExists bool
		wantErr               error
	}{
		{
			name: "refresh tokenが有効な場合",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(&entity.RefreshToken{
						ID:        parsedUUID,
						UserID:    parsedUUID,
						TokenHash: tokenHash,
						ExpiresAt: now.Add(time.Hour),
						RevokedAt: nil,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil)
			},
			wantAccessTokenExists: true,
			wantErr:               nil,
		},
		{
			name: "refresh token検索に失敗する",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(nil, findErr)
			},
			wantAccessTokenExists: false,
			wantErr:               findErr,
		},
		{
			name: "refresh tokenが存在しない",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(nil, nil)
			},
			wantAccessTokenExists: false,
			wantErr:               ErrRefreshTokenNotFound,
		},
		{
			name: "refresh tokenがrevoked",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				revokedAt := now
				authRepo.EXPECT().FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(&entity.RefreshToken{
						ID:        parsedUUID,
						UserID:    parsedUUID,
						TokenHash: tokenHash,
						ExpiresAt: now.Add(time.Hour),
						RevokedAt: &revokedAt,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil)
			},
			wantAccessTokenExists: false,
			wantErr:               ErrRefreshTokenRevoked,
		},
		{
			name: "refresh tokenが期限切れ",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(&entity.RefreshToken{
						ID:        parsedUUID,
						UserID:    parsedUUID,
						TokenHash: tokenHash,
						ExpiresAt: now.Add(-time.Hour),
						RevokedAt: nil,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil)
			},
			wantAccessTokenExists: false,
			wantErr:               ErrRefreshTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authRepo := mock.NewMockAuthRepository(ctrl)
			tt.mockFunc(authRepo)

			idTokenValidator := authmock.NewMockIDTokenValidator(ctrl)
			tx := mock.NewMockTransaction(ctrl)

			uc := NewAuthUsecase(
				authRepo,
				&mock.UserRepositoryMock{},
				&mock.CharacterRepositoryMock{},
				jwtService,
				idTokenValidator,
				tx,
			)

			got, err := uc.RefreshAccessToken(ctx, refreshToken)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantAccessTokenExists, got.AccessToken != "")
		})
	}
}

func Test_AuthUsecase_LoginWithGoogle(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	jwtService := auth.NewJWTService("test-secret", time.Hour, "test-issuer")

	tests := []struct {
		name     string
		mockFunc func(ctrl *gomock.Controller) (
			authRepo *mock.MockAuthRepository,
			idTokenValidator *authmock.MockIDTokenValidator,
			userRepo *mock.UserRepositoryMock,
			characterRepo *mock.CharacterRepositoryMock,
			tx *mock.MockTransaction,
		)
		wantAccessTokenExists  bool
		wantRefreshTokenExists bool
		wantErr                error
	}{
		{
			name: "既存ユーザーでログイン成功",
			mockFunc: func(ctrl *gomock.Controller) (
				authRepo *mock.MockAuthRepository,
				idTokenValidator *authmock.MockIDTokenValidator,
				userRepo *mock.UserRepositoryMock,
				characterRepo *mock.CharacterRepositoryMock,
				tx *mock.MockTransaction,
			) {
				authRepo = mock.NewMockAuthRepository(ctrl)
				userRepo = &mock.UserRepositoryMock{}
				characterRepo = &mock.CharacterRepositoryMock{}
				idTokenValidator = authmock.NewMockIDTokenValidator(ctrl)
				tx = mock.NewMockTransaction(ctrl)

				idTokenValidator.EXPECT().
					Validate(gomock.Any(), "valid-token").
					Return(&domainauth.IDTokenPayload{
						Subject: "google-sub",
					}, nil)

				authRepo.EXPECT().
					FindByProviderAndProviderUserID(
						gomock.Any(),
						entity.AuthProviderGoogle,
						"google-sub",
					).
					Return(&entity.AuthIdentity{
						UserID: parsedUUID,
					}, nil)

				authRepo.EXPECT().
					CreateRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil)

				return authRepo, idTokenValidator, userRepo, characterRepo, tx
			},
			wantAccessTokenExists:  true,
			wantRefreshTokenExists: true,
			wantErr:                nil,
		},
		{
			name: "新規ユーザーでログイン成功",
			mockFunc: func(ctrl *gomock.Controller) (
				authRepo *mock.MockAuthRepository,
				idTokenValidator *authmock.MockIDTokenValidator,
				userRepo *mock.UserRepositoryMock,
				characterRepo *mock.CharacterRepositoryMock,
				tx *mock.MockTransaction,
			) {
				authRepo = mock.NewMockAuthRepository(ctrl)
				tx = mock.NewMockTransaction(ctrl)

				userRepo = &mock.UserRepositoryMock{
					CreateFunc: func(ctx context.Context, user entity.User) error {
						assert.Equal(t, "new user", user.Nickname)
						return nil
					},
				}

				characterRepo = &mock.CharacterRepositoryMock{
					CreateFunc: func(ctx context.Context, character entity.Character) error {
						assert.Equal(t, entity.SupporterTypeDog, character.SupporterType)
						return nil
					},
				}

				idTokenValidator = authmock.NewMockIDTokenValidator(ctrl)

				idTokenValidator.EXPECT().
					Validate(gomock.Any(), "valid-token").
					Return(&domainauth.IDTokenPayload{
						Subject: "google-sub",
					}, nil)

				authRepo.EXPECT().
					FindByProviderAndProviderUserID(
						gomock.Any(),
						entity.AuthProviderGoogle,
						"google-sub",
					).
					Return(nil, nil)

				authRepo.EXPECT().
					CreateAuthIdentity(gomock.Any(), gomock.Any()).
					Return(nil)

				authRepo.EXPECT().
					CreateRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil)

				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				return authRepo, idTokenValidator, userRepo, characterRepo, tx
			},
			wantAccessTokenExists:  true,
			wantRefreshTokenExists: true,
			wantErr:                nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authRepo := mock.NewMockAuthRepository(ctrl)
			idTokenValidator := authmock.NewMockIDTokenValidator(ctrl)
			authRepo, idTokenValidator, userRepo, characterRepo, tx := tt.mockFunc(ctrl)
			uc := NewAuthUsecase(
				authRepo,
				userRepo,
				characterRepo,
				jwtService,
				idTokenValidator,
				tx,
			)
			got, err := uc.LoginWithGoogle(ctx, "valid-token")
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantAccessTokenExists, got.AccessToken != "")
			assert.Equal(t, tt.wantRefreshTokenExists, got.RefreshToken != "")
		})
	}
}

func Test_AuthUsecase_Logout(t *testing.T) {
	ctx := context.Background()

	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	refreshToken := "test-refresh-token"
	tokenHash := hashRefreshToken(refreshToken)
	now := time.Now().UTC().Truncate(time.Millisecond)

	findErr := errors.New("failed to find refresh token")
	revokeErr := errors.New("failed to revoke refresh token")

	jwtService := auth.NewJWTService("test-secret", time.Hour, "test-issur")

	tests := []struct {
		name     string
		mockFunc func(*mock.MockAuthRepository)
		wantErr  error
	}{
		{
			name: "refresh tokenが見つかりlogout成功",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().
					FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(&entity.RefreshToken{
						ID:        id,
						UserID:    id,
						TokenHash: tokenHash,
						ExpiresAt: now.Add(time.Hour),
						RevokedAt: nil,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil)
				authRepo.EXPECT().
					RevokeRefreshToken(gomock.Any(), id, gomock.Any()).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "refresh tokenが見つからない場合も成功扱い",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().
					FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(nil, nil)
			},
			wantErr: nil,
		},
		{
			name: "refresh token検索に失敗する",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().
					FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(nil, findErr)
			},
			wantErr: findErr,
		},
		{
			name: "refresh tokenのrevokeに失敗する",
			mockFunc: func(authRepo *mock.MockAuthRepository) {
				authRepo.EXPECT().
					FindRefreshTokenByHash(gomock.Any(), tokenHash).
					Return(&entity.RefreshToken{
						ID:        id,
						UserID:    id,
						TokenHash: tokenHash,
						ExpiresAt: now.Add(time.Hour),
						RevokedAt: nil,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil)

				authRepo.EXPECT().
					RevokeRefreshToken(gomock.Any(), id, gomock.Any()).
					Return(revokeErr)
			},
			wantErr: revokeErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authRepo := mock.NewMockAuthRepository(ctrl)
			tt.mockFunc(authRepo)
			idTokenValidator := authmock.NewMockIDTokenValidator(ctrl)
			tx := mock.NewMockTransaction(ctrl)
			uc := NewAuthUsecase(
				authRepo,
				&mock.UserRepositoryMock{},
				&mock.CharacterRepositoryMock{},
				jwtService,
				idTokenValidator,
				tx,
			)

			err := uc.Logout(ctx, refreshToken)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}

}
