package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	mock "github.com/chocoko/character-support-walk-backend/internal/test/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_WalkOptionUsecaes_List(t *testing.T) {
	ctx := context.Background()
	findErr := errors.New("failed to find walker")

	walkOptions := []entity.WalkOption{
		{
			ID:       "1",
			Category: "free",
			Title:    "ふらっとそこらへん",
		},
		{
			ID:       "2",
			Category: "minutes",
			Title:    "5分だけ歩く",
		},
	}

	tests := []struct {
		name     string
		mockFunc func(repo *mock.MockWalkOptionRepository)
		want     []entity.WalkOption
		wantErr  error
	}{
		{
			name: "WalkOptionの取得",
			mockFunc: func(repo *mock.MockWalkOptionRepository) {
				repo.EXPECT().SelectOption(gomock.Any()).Return(walkOptions, nil)
			},
			want:    walkOptions,
			wantErr: nil,
		},
		{
			name: "WalkOptionの取得失敗",
			mockFunc: func(repo *mock.MockWalkOptionRepository) {
				repo.EXPECT().SelectOption(gomock.Any()).Return(nil, findErr)
			},
			want:    nil,
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockWalkOptionRepository(ctrl)
			tt.mockFunc(repo)
			uc := NewWalkOptionUsecase(repo)
			got, err := uc.List(ctx)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
