package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/chocoko/character-support-walk-backend/internal/test/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewAuthRepository(t *testing.T) {
	db, _ := testutils.NewMockDB(t)
	repo := NewAuthRepository(db)
	assert.Equal(t, db, repo.db)
}

func Test_AuthIdentityRepository_FindByProviderAndProviderUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found auth user")

	tests := []struct {
		name           string
		provider       entity.AuthProvider
		providerUserID string
		mockFunc       func(sqlmock.Sqlmock)
		want           *entity.AuthIdentity
		wantErr        error
	}{
		{
			name:           "対象のAuthUserが見つかる",
			provider:       entity.AuthProviderGoogle,
			providerUserID: "test_provider_user_id",
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "provider", "provider_user_id", "created_at", "updated_at"})
				rows.AddRow(id, id, entity.AuthProviderGoogle, "test_provider_user_id", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
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
				`)).WithArgs("test_provider_user_id", entity.AuthProviderGoogle).WillReturnRows(rows)
			},
			want: &entity.AuthIdentity{
				ID:             id,
				UserID:         id,
				Provider:       entity.AuthProviderGoogle,
				ProviderUserID: "test_provider_user_id",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			wantErr: nil,
		},
		{
			name:           "対象のAuthUser取得時DBエラー",
			provider:       entity.AuthProviderGoogle,
			providerUserID: "test_provider_user_id",
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "provider", "provider_user_id", "created_at", "updated_at"})
				rows.AddRow(id, id, entity.AuthProviderGoogle, "test_provider_user_id", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
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
				`)).WithArgs("test_provider_user_id", entity.AuthProviderGoogle).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
		{
			name:           "対象のAuthUserが見つからない",
			provider:       entity.AuthProviderGoogle,
			providerUserID: "test_provider_user_id",
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "provider", "provider_user_id", "created_at", "updated_at"})
				rows.AddRow(id, id, entity.AuthProviderGoogle, "test_provider_user_id", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
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
				`)).WithArgs("test_provider_user_id", entity.AuthProviderGoogle).WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: repository.ErrAuthIdentityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewAuthRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindByProviderAndProviderUserID(ctx, tt.provider, tt.providerUserID)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_AuthIdentityRepository_CreateAuthIdentity(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	insertErr := errors.New("failed to insert data")

	authIdentity := &entity.AuthIdentity{
		ID:             id,
		UserID:         id,
		Provider:       entity.AuthProviderGoogle,
		ProviderUserID: "test_provider_user_id",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	tests := []struct {
		name         string
		authIdentity *entity.AuthIdentity
		mockFunc     func(sqlmock.Sqlmock)
		wantErr      error
	}{
		{
			name:         "AuthUserの登録が成功したとき",
			authIdentity: authIdentity,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				    INSERT INTO auth_identities (
						id,
						user_id,
						provider,
						provider_user_id,
						created_at,
						updated_at 
					)
					VALUES($1, $2, $3, $4, $5, $6)
				`)).WithArgs(
					authIdentity.ID,
					authIdentity.UserID,
					authIdentity.Provider,
					authIdentity.ProviderUserID,
					authIdentity.CreatedAt,
					authIdentity.UpdatedAt,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name:         "AuthUserの登録が失敗したとき",
			authIdentity: authIdentity,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  INSERT INTO auth_identities (
					  id,
					  user_id,
					  provider,
				      provider_user_id,
					  created_at,
					  updated_at 
				  )
				  VALUES($1, $2, $3, $4, $5, $6)
				`)).WithArgs(
					authIdentity.ID,
					authIdentity.UserID,
					authIdentity.Provider,
					authIdentity.ProviderUserID,
					authIdentity.CreatedAt,
					authIdentity.UpdatedAt,
				).WillReturnError(insertErr)
			},
			wantErr: insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewAuthRepository(db)
			tt.mockFunc(mock)
			err := repo.CreateAuthIdentity(ctx, *tt.authIdentity)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_AuthIdentityRepository_FindRefreshTokenByHash(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	threeHoursAgo := now.Add(-3 * time.Hour)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found auth user")

	tests := []struct {
		name      string
		tokenHash string
		mockFunc  func(sqlmock.Sqlmock)
		want      *entity.RefreshToken
		wantErr   error
	}{
		{
			name:      "対象のUserが見つかる",
			tokenHash: "test_refresh_token",
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "expires_at", "revoked_at", "created_at", "updated_at"})
				rows.AddRow(id, id, "test_refresh_token", threeHoursAgo, nil, now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
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
				`)).WithArgs("test_refresh_token").WillReturnRows(rows)
			},
			want: &entity.RefreshToken{
				ID:        id,
				UserID:    id,
				TokenHash: "test_refresh_token",
				ExpiresAt: threeHoursAgo,
				RevokedAt: nil,
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: nil,
		},
		{
			name:      "DBエラーが発生した時",
			tokenHash: "test_refresh_token",
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "expires_at", "revoked_at", "created_at", "updated_at"})
				rows.AddRow(id, id, "test_refresh_token", threeHoursAgo, nil, now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
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
				`)).WithArgs("test_refresh_token").WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
		{
			name:      "対象のUserが見つからない",
			tokenHash: "test_refresh_token",
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "expires_at", "revoked_at", "created_at", "updated_at"})
				rows.AddRow(id, id, "test_refresh_token", threeHoursAgo, nil, now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
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
				`)).WithArgs("test_refresh_token").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: repository.ErrRefreshTokenNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewAuthRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindRefreshTokenByHash(ctx, tt.tokenHash)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_AuthIndetityRepository_CreateRefreshToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	threeHoursAgo := now.Add(-3 * time.Hour)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	insertErr := errors.New("failed to insert data")

	refreshToken := &entity.RefreshToken{
		ID:        id,
		UserID:    id,
		TokenHash: "test_refresh_token",
		ExpiresAt: threeHoursAgo,
		RevokedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
	tests := []struct {
		name         string
		refreshToken *entity.RefreshToken
		mockFunc     func(sqlmock.Sqlmock)
		wantErr      error
	}{
		{
			name:         "UserのrefreshTokenが登録に成功したとき",
			refreshToken: refreshToken,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
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
				`)).WithArgs(
					refreshToken.ID,
					refreshToken.UserID,
					refreshToken.TokenHash,
					refreshToken.ExpiresAt,
					refreshToken.RevokedAt,
					refreshToken.CreatedAt,
					refreshToken.UpdatedAt,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name:         "UserのrefreshTokenが登録に失敗したとき",
			refreshToken: refreshToken,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
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
				`)).WithArgs(
					refreshToken.ID,
					refreshToken.UserID,
					refreshToken.TokenHash,
					refreshToken.ExpiresAt,
					refreshToken.RevokedAt,
					refreshToken.CreatedAt,
					refreshToken.UpdatedAt,
				).WillReturnError(insertErr)
			},
			wantErr: insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewAuthRepository(db)
			tt.mockFunc(mock)
			err := repo.CreateRefreshToken(ctx, *tt.refreshToken)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_AuthIdentityRepository_RevokeRefreshToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	threeHoursAgo := now.Add(-3 * time.Hour)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	updateErr := errors.New("failed to update data")

	tests := []struct {
		name           string
		refreshTokenID uuid.UUID
		revokedAt      time.Time
		mockFunc       func(sqlmock.Sqlmock)
		wantErr        error
	}{
		{
			name:           "RefreshTokenの更新が成功したとき",
			refreshTokenID: id,
			revokedAt:      threeHoursAgo,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  UPDATE refresh_tokens
				  SET
					revoked_at = $1,
					updated_at = $2
				  WHERE id = $3
				`)).WithArgs(
					threeHoursAgo,
					threeHoursAgo,
					id,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name:           "RefreshTokenの更新が失敗したとき",
			refreshTokenID: id,
			revokedAt:      threeHoursAgo,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  UPDATE refresh_tokens
				  SET
					revoked_at = $1,
					updated_at = $2
				  WHERE id = $3
				`)).WithArgs(
					threeHoursAgo,
					threeHoursAgo,
					id,
				).WillReturnError(updateErr)
			},
			wantErr: updateErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewAuthRepository(db)
			tt.mockFunc(mock)
			err := repo.RevokeRefreshToken(ctx, tt.refreshTokenID, tt.revokedAt)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
