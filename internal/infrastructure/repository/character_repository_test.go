package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/test/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewCharacterRepository(t *testing.T) {
	db, _ := testutils.NewMockDB(t)
	repo := NewCharacterRepository(db)
	assert.Equal(t, db, repo.db)
}

func Test_CharacterRepository_FindByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found character")

	character := entity.Character{
		ID:            id,
		UserID:        id,
		SupporterType: entity.SupporterTypeDog,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name     string
		userID   uuid.UUID
		want     *entity.Character
		mockFunc func(sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name:   "指定したUserIDに該当するデータの取得",
			userID: id,
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "supporter_type", "created_at", "updated_at"})
				rows.AddRow(id, id, entity.SupporterTypeDog, now, now)
				s.ExpectQuery(regexp.QuoteMeta(`
				SElECT
					id,
					user_id,
					supporter_type,
					created_at,
					updated_at
				FROM characters
				WHERE user_id = $1
				`)).WithArgs(id).WillReturnRows(rows)
			},
			want:    &character,
			wantErr: nil,
		},
		{
			name:   "指定したUserIDに該当するデータの取得の失敗",
			userID: id,
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "supporter_type", "created_at", "updated_at"})
				rows.AddRow(id, id, entity.SupporterTypeDog, now, now)
				s.ExpectQuery(regexp.QuoteMeta(`
				SElECT
					id,
					user_id,
					supporter_type,
					created_at,
					updated_at
				FROM characters
				WHERE user_id = $1
				`)).WithArgs(id).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewCharacterRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindByUserID(ctx, tt.userID)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_CharacterRepository_UpdateByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found character")

	character := entity.Character{
		UserID:        id,
		SupporterType: entity.SupporterTypeDog,
		UpdatedAt:     now,
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		supporterType entity.SupporterType
		updatedAt     time.Time
		want          *entity.Character
		mockFunc      func(sqlmock.Sqlmock)
		wantErr       error
	}{
		{
			name:          "指定したUserIDに該当するデータの更新",
			userID:        id,
			supporterType: entity.SupporterTypeDog,
			updatedAt:     now,
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id", "supporter_type", "updated_at"})
				rows.AddRow(id, entity.SupporterTypeDog, now)
				s.ExpectQuery(regexp.QuoteMeta(`
				UPDATE characters
				SET
					supporter_type = $1,
					updated_at = $2
				WHERE user_id = $3
				RETURNING user_id, supporter_type, updated_at;
				`)).WithArgs(entity.SupporterTypeDog, now, id).WillReturnRows(rows)
			},
			want:    &character,
			wantErr: nil,
		},
		{
			name:          "指定したUserIDに該当するデータの更新の失敗",
			userID:        id,
			supporterType: entity.SupporterTypeDog,
			updatedAt:     now,
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id", "supporter_type", "updated_at"})
				rows.AddRow(id, entity.SupporterTypeDog, now)
				s.ExpectQuery(regexp.QuoteMeta(`
				UPDATE characters
				SET
					supporter_type = $1,
					updated_at = $2
				WHERE user_id = $3
				RETURNING user_id, supporter_type, updated_at;
				`)).WithArgs(entity.SupporterTypeDog, now, id).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewCharacterRepository(db)
			tt.mockFunc(mock)
			got, err := repo.UpdateByUserID(ctx, tt.userID, tt.supporterType, tt.updatedAt)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_CharacterRepository_Create(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	insertErr := errors.New("failed to insert character")

	character := entity.Character{
		ID:            id,
		UserID:        id,
		SupporterType: entity.SupporterTypeDog,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tests := []struct {
		name      string
		character entity.Character
		mockFunc  func(sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name:      "Characterデータの登録",
			character: character,
			mockFunc: func(s sqlmock.Sqlmock) {
				s.ExpectExec(regexp.QuoteMeta(`
				INSERT INTO characters (
					id,
					user_id,
					supporter_type,
					created_at,
					updated_at 
				)
				VALUES($1, $2, $3, $4, $5)
				`)).WithArgs(character.ID, character.UserID, character.SupporterType, character.CreatedAt, character.UpdatedAt).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name:      "Characterデータの登録の失敗",
			character: character,
			mockFunc: func(s sqlmock.Sqlmock) {
				sqlmock.NewRows([]string{"id", "user_id", "supporter_type", "created_at", "updated_at"})
				s.ExpectExec(regexp.QuoteMeta(`
				INSERT INTO characters (
					id,
					user_id,
					supporter_type,
					created_at,
					updated_at 
				)
				VALUES($1, $2, $3, $4, $5)
				`)).WithArgs(id, id, entity.SupporterTypeDog, now, now).WillReturnError(insertErr)
			},
			wantErr: insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewCharacterRepository(db)
			tt.mockFunc(mock)
			err := repo.Create(ctx, tt.character)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
