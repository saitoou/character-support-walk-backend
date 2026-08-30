package repository

import (
	"context"
	"testing"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func Test_WalkOptionRepositoryIntegration_SelectOption(t *testing.T) {

	walkOption := []entity.WalkOption{
		{
			ID:       "1",
			Category: "free",
			Title:    "ふらっとそこらへん",
		},
		{
			ID:       "2",
			Category: "minutes",
			Title:    "5分だけ外に出る",
		},
		{
			ID:       "3",
			Category: "destination",
			Title:    "コンビニまで",
		},
	}

	tests := []struct {
		name     string
		setupCtx func(t *testing.T) context.Context
		want     []entity.WalkOption
		wantErr  error
	}{
		{
			name: "WalkOptionの取得を行う",
			setupCtx: func(t *testing.T) context.Context {
				t.Helper()
				return t.Context()
			},
			want:    walkOption,
			wantErr: nil,
		},
		{
			name: "WalkOptionの取得に失敗する",
			setupCtx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			want:    nil,
			wantErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := beginTestTx(t)
			repo := NewWalkOptionRepository(tx)
			ctx := tt.setupCtx(t)
			got, err := repo.SelectOption(ctx)

			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
