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

func Test_NewUserRepository(t *testing.T) {
	db, _ := testutils.NewMockDB(t)
	repo := NewUserRepository(db)
	assert.Equal(t, db, repo.db)
}

func Test_UserRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found user")

	tests := []struct {
		name     string
		userID   uuid.UUID
		mockFunc func(sqlmock.Sqlmock)
		want     *entity.User
		wantErr  error
	}{
		{
			name:   "対象のUserが見つかる",
			userID: id,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nickname", "created_at", "updated_at"})
				rows.AddRow(id, "test user", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SELECT
						id,
						nickname,
						created_at,
						updated_at
					FROM users
					WHERE id = $1 AND deleted_at IS NULL;
				`)).WithArgs(id).WillReturnRows(rows)
			},
			want: &entity.User{
				ID:        id,
				Nickname:  "test user",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: nil,
		},
		{
			name:   "対象のUserが見つからない",
			userID: id,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nickname", "created_at", "updated_at"})
				rows.AddRow(id, "test user", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SELECT
						id,
						nickname,
						created_at,
						updated_at
					FROM users
					WHERE id = $1 AND deleted_at IS NULL;
				`)).WithArgs(id).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewUserRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindByID(ctx, tt.userID)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_UserRepository_Create(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	insertErr := errors.New("failed to insert data")

	user := &entity.User{
		ID:        id,
		Nickname:  "test user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name     string
		user     *entity.User
		mockFunc func(sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "Userの登録が成功したとき",
			user: user,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				    INSERT INTO users (
						id,
						nickname,
						created_at,
						updated_at 
					)
					VALUES($1, $2, $3, $4)
				`)).WithArgs(
					user.ID,
					user.Nickname,
					user.CreatedAt,
					user.UpdatedAt,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name: "Userの登録が失敗したとき",
			user: user,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  INSERT INTO users (
						id,
						nickname,
						created_at,
						updated_at 
					)
					VALUES($1, $2, $3, $4)
				`)).WithArgs(
					user.ID,
					user.Nickname,
					user.CreatedAt,
					user.UpdatedAt,
				).WillReturnError(insertErr)
			},
			wantErr: insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewUserRepository(db)
			tt.mockFunc(mock)
			err := repo.Create(ctx, *tt.user)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_UserRepository_Update(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	updateErr := errors.New("failed to update user")

	user := entity.User{
		ID:        id,
		Nickname:  "test user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		want      *entity.User
		userID    uuid.UUID
		nickname  string
		updatedAt time.Time
		mockFunc  func(sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name:      "WalkingのStatus更新が成功したとき",
			want:      &user,
			userID:    id,
			nickname:  "test user",
			updatedAt: now,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nickname", "created_at", "updated_at"})
				rows.AddRow(id, "test user", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    UPDATE users
					SET
						nickname = $1,
						updated_at = $2
					WHERE id = $3 AND deleted_at IS NULL
					RETURNING id, nickname, created_at, updated_at;
				`)).WithArgs(
					"test user",
					now,
					id,
				).WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name:      "WalkingのStatus更新が失敗したとき",
			want:      nil,
			userID:    id,
			nickname:  "test user",
			updatedAt: now,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nickname", "created_at", "updated_at"})
				rows.AddRow(id, "test user", now, now)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    UPDATE users
					SET
						nickname = $1,
						updated_at = $2
					WHERE id = $3 AND deleted_at IS NULL
					RETURNING id, nickname, created_at, updated_at;
				`)).WithArgs(
					"test user",
					now,
					id,
				).WillReturnError(updateErr)
			},
			wantErr: updateErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewUserRepository(db)
			tt.mockFunc(mock)
			got, err := repo.Update(ctx, tt.userID, tt.nickname, tt.updatedAt)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}
