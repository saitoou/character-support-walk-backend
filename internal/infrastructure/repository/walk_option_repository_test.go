package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/test/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_NewWalkOptionRepository(t *testing.T) {
	db, _ := testutils.NewMockDB(t)
	repo := NewWalkRepository(db)
	assert.Equal(t, db, repo.db)
}

func Test_WalkOptionRepository_SelectOption(t *testing.T) {
	ctx := context.Background()
	selectErr := errors.New("failed to insert data")

	walkOption := []entity.WalkOption{
		{
			ID:       "1",
			Category: "Walking",
			Title:    "今日のおさんぽ",
		},
		{
			ID:       "2",
			Category: "Complete",
			Title:    "今日のおさんぽ",
		},
	}

	tests := []struct {
		name     string
		mockFunc func(sqlmock.Sqlmock)
		want     []entity.WalkOption
		wantErr  error
	}{
		{
			name: "WalkOptionの取得を行う",
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "category", "title"})
				rows.AddRow("1", "Walking", "今日のおさんぽ")
				rows.AddRow("2", "Complete", "今日のおさんぽ")
				s.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, category, title FROM walk_options ORDER BY id ASC;
				`)).WillReturnRows(rows)
			},
			want:    walkOption,
			wantErr: nil,
		},
		{
			name: "WalkOptionの取得に失敗する",
			mockFunc: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, category, title FROM walk_options ORDER BY id ASC;
				`)).WillReturnError(selectErr)
			},
			want:    nil,
			wantErr: selectErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := testutils.NewMockDB(t)
			repo := NewWalkOptionRepository(db)
			tt.mockFunc(mock)
			got, err := repo.SelectOption(ctx)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
