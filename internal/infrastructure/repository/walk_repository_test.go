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

func Test_NewWalkRespository(t *testing.T) {
	db, _ := testutils.NewMockDB(t)
	repo := NewWalkRepository(db)
	assert.Equal(t, db, repo.db)
}

func Test_WalkRepository_Create(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	insertErr := errors.New("failed to insert data")

	walk := &entity.Walk{
		ID:           id,
		UserID:       id,
		WalkOptionID: 1,
		Status:       entity.WalkStatusWalking,
		StartedAt:    now,
		FinishedAt:   nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name     string
		walk     *entity.Walk
		mockFunc func(sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "Walkの登録が成功したとき",
			walk: walk,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  INSERT INTO walks (
				      id,
					  user_id,
					  walk_option_id,
					  status,
					  started_at,
					  finished_at,
					  created_at,
					  updated_at 
				  ) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
				`)).WithArgs(
					walk.ID,
					walk.UserID,
					walk.WalkOptionID,
					walk.Status,
					walk.StartedAt,
					walk.FinishedAt,
					walk.CreatedAt,
					walk.UpdatedAt,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name: "Walkの登録が失敗したとき",
			walk: walk,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  INSERT INTO walks (
				      id,
					  user_id,
					  walk_option_id,
					  status,
					  started_at,
					  finished_at,
					  created_at,
					  updated_at 
				  ) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
				`)).WithArgs(
					walk.ID,
					walk.UserID,
					walk.WalkOptionID,
					walk.Status,
					walk.StartedAt,
					walk.FinishedAt,
					walk.CreatedAt,
					walk.UpdatedAt,
				).WillReturnError(insertErr)
			},
			wantErr: insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewWalkRepository(db)
			tt.mockFunc(mock)
			err := repo.Create(ctx, *tt.walk)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkRepository_FindByUserID(t *testing.T) {
	ctx := context.Background()
	recordDate := time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)
	threeDaysAgo := recordDate.Add(-3 * time.Hour)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found walk")

	tests := []struct {
		name     string
		id       uuid.UUID
		mockFunc func(sqlmock.Sqlmock)
		want     []entity.Walk
		wantErr  error
	}{
		{
			name: "対象のWalkが見つかる",
			id:   id,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "walk_option_id", "status", "started_at", "finished_at", "created_at", "updated_at"})
				rows.AddRow(id, id, 1, entity.WalkStatusWalking, recordDate, nil, recordDate, recordDate)
				rows.AddRow(id, id, 1, entity.WalkStatusCompleted, threeDaysAgo, threeDaysAgo, threeDaysAgo, threeDaysAgo)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SElECT
						id,
						user_id,
						walk_option_id,
						status,
						started_at,
						finished_at,
						created_at,
						updated_at
					FROM walks
					WHERE user_id = $1
					ORDER BY created_at DESC
				`)).WithArgs(id).WillReturnRows(rows)
			},
			want: []entity.Walk{
				{
					ID:           id,
					UserID:       id,
					WalkOptionID: 1,
					Status:       entity.WalkStatusWalking,
					StartedAt:    recordDate,
					FinishedAt:   nil,
					CreatedAt:    recordDate,
					UpdatedAt:    recordDate,
				},
				{
					ID:           id,
					UserID:       id,
					WalkOptionID: 1,
					Status:       entity.WalkStatusCompleted,
					StartedAt:    threeDaysAgo,
					FinishedAt:   &threeDaysAgo,
					CreatedAt:    threeDaysAgo,
					UpdatedAt:    threeDaysAgo,
				},
			},
			wantErr: nil,
		},
		{
			name: "対象のWalkが見つからない",
			id:   id,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SElECT
						id,
						user_id,
						walk_option_id,
						status,
						started_at,
						finished_at,
						created_at,
						updated_at
					FROM walks
					WHERE user_id = $1
					ORDER BY created_at DESC
				`)).WithArgs(id).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewWalkRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindByUserID(ctx, tt.id)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkRepository_FindByUserIDAndID(t *testing.T) {
	ctx := context.Background()
	recordDate := time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found walk")

	tests := []struct {
		name     string
		id       uuid.UUID
		userID   uuid.UUID
		mockFunc func(sqlmock.Sqlmock)
		want     *entity.Walk
		wantErr  error
	}{
		{
			name:   "対象UserIDとWalkIDに該当するWalkが見つかる",
			id:     id,
			userID: id,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "walk_option_id", "status", "started_at", "finished_at", "created_at", "updated_at"})
				rows.AddRow(id, id, 1, entity.WalkStatusWalking, recordDate, nil, recordDate, recordDate)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SELECT
						id,
						user_id,
						walk_option_id,
						status,
						started_at,
						finished_at,
						created_at,
						updated_at
					FROM walks
					WHERE user_id = $1
						AND ID = $2
				`)).WithArgs(id, id).WillReturnRows(rows)
			},
			want: &entity.Walk{
				ID:           id,
				UserID:       id,
				WalkOptionID: 1,
				Status:       entity.WalkStatusWalking,
				StartedAt:    recordDate,
				FinishedAt:   nil,
				CreatedAt:    recordDate,
				UpdatedAt:    recordDate,
			},
			wantErr: nil,
		},
		{
			name:   "対象のUserIDとWalkIDに該当するWalkが見つからない",
			id:     id,
			userID: id,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SELECT
						id,
						user_id,
						walk_option_id,
						status,
						started_at,
						finished_at,
						created_at,
						updated_at
					FROM walks
					WHERE user_id = $1
						AND ID = $2
				`)).WithArgs(id, id).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewWalkRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindByUserIDAndID(ctx, tt.userID, tt.id)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkRespository_FindTodayWalkByUserID(t *testing.T) {
	ctx := context.Background()
	recordDate := time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)
	threeHoursAgo := recordDate.Add(-3 * time.Hour)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	findErr := errors.New("not found walk")

	tests := []struct {
		name          string
		userID        uuid.UUID
		todayStart    time.Time
		tomorrowStart time.Time
		mockFunc      func(sqlmock.Sqlmock)
		want          *entity.Walk
		wantErr       error
	}{
		{
			name:          "対象のUserIDに該当するWalkが見つかる",
			userID:        id,
			todayStart:    recordDate,
			tomorrowStart: threeHoursAgo,
			mockFunc: func(sm sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "walk_option_id", "status", "started_at", "finished_at", "created_at", "updated_at"})
				rows.AddRow(id, id, 1, entity.WalkStatusWalking, recordDate, nil, recordDate, recordDate)
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SElECT
						id,
						user_id,
						walk_option_id,
						status,
						started_at,
						finished_at,
						created_at,
						updated_at
					FROM walks
					WHERE user_id = $1
						AND created_at >= $2
						AND created_at < $3
					ORDER BY
						CASE
						WHEN status = 'completed' THEN 1
						WHEN status = 'walking' THEN 2
						WHEN status = 'canceled' THEN 3
						ELSE 4
						END ASC,
						created_at DESC
					LIMIT 1
				`)).WithArgs(id, recordDate, threeHoursAgo).WillReturnRows(rows)
			},
			want: &entity.Walk{
				ID:           id,
				UserID:       id,
				WalkOptionID: 1,
				Status:       entity.WalkStatusWalking,
				StartedAt:    recordDate,
				FinishedAt:   nil,
				CreatedAt:    recordDate,
				UpdatedAt:    recordDate,
			},
			wantErr: nil,
		},
		{
			name:          "対象のUserIDに該当するWalkが見つからない",
			userID:        id,
			todayStart:    recordDate,
			tomorrowStart: threeHoursAgo,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectQuery(regexp.QuoteMeta(`
				    SElECT
						id,
						user_id,
						walk_option_id,
						status,
						started_at,
						finished_at,
						created_at,
						updated_at
					FROM walks
					WHERE user_id = $1
						AND created_at >= $2
						AND created_at < $3
					ORDER BY
						CASE
						WHEN status = 'completed' THEN 1
						WHEN status = 'walking' THEN 2
						WHEN status = 'canceled' THEN 3
						ELSE 4
						END ASC,
						created_at DESC
					LIMIT 1
				`)).WithArgs(id, recordDate, threeHoursAgo).WillReturnError(findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewWalkRepository(db)
			tt.mockFunc(mock)
			got, err := repo.FindTodayWalkByUserID(ctx, tt.userID, tt.todayStart, tt.tomorrowStart)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkRespository_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	id := uuid.MustParse("019e1cd3-8194-7a36-816b-2f38206ca52d")
	insertErr := errors.New("failed to insert data")
	completeStatus := entity.WalkStatusCompleted

	walk := entity.Walk{
		ID:           id,
		UserID:       id,
		WalkOptionID: 1,
		Status:       entity.WalkStatusCompleted,
		StartedAt:    now,
		FinishedAt:   &now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name     string
		walk     entity.Walk
		mockFunc func(sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "WalkingのStatus更新が成功したとき",
			walk: walk,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  UPDATE walks
				  SET
					status = $1,
					finished_at = $2,
					updated_at = $3
				  WHERE id = $4
				    AND user_id = $5
					AND status = 'walking'
				`)).WithArgs(
					completeStatus,
					now,
					now,
					id,
					id,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name: "WalkingのStatus更新が失敗したとき",
			walk: walk,
			mockFunc: func(sm sqlmock.Sqlmock) {
				sm.ExpectExec(regexp.QuoteMeta(`
				  UPDATE walks
				  SET
					status = $1,
					finished_at = $2,
					updated_at = $3
				  WHERE id = $4
				    AND user_id = $5
					AND status = 'walking'
				`)).WithArgs(
					completeStatus,
					now,
					now,
					id,
					id,
				).WillReturnError(insertErr)
			},
			wantErr: insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewWalkRepository(db)
			tt.mockFunc(mock)
			err := repo.UpdateStatus(ctx, tt.walk)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
