package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	mock "github.com/chocoko/character-support-walk-backend/internal/test/mock/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_HomeUsecase_GetHome(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	findErr := errors.New("failed to find home")
	now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name     string
		UserID   string
		mockFunc func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock, *mock.WalkRepositoryMock)
		want     HomeOutput
		wantErr  error
	}{
		{
			name:   "Home一覧を取得する",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock, *mock.WalkRepositoryMock) {
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        parsedUUID,
								Nickname:  "test user",
								CreatedAt: now,
								UpdatedAt: now,
							}, nil
						},
					}, &mock.CharacterRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
							return &entity.Character{
								ID:            parsedUUID,
								UserID:        parsedUUID,
								SupporterType: entity.SupporterTypeDog,
								CreatedAt:     now,
								UpdatedAt:     now,
							}, nil
						},
					}, &mock.WalkRepositoryMock{
						FindTodayWalkByUserIDFunc: func(ctx context.Context, userID uuid.UUID, todayStart, tomorrowStart time.Time) (*entity.Walk, error) {
							return &entity.Walk{
								ID:           parsedUUID,
								UserID:       parsedUUID,
								WalkOptionID: 1,
								Status:       entity.WalkStatusWalking,
								StartedAt:    now,
								FinishedAt:   nil,
								CreatedAt:    now,
								UpdatedAt:    now,
							}, nil
						},
					}
			},
			want: HomeOutput{
				Walker: HomeWalkerOutput{
					Nickname: "test user",
				},
				Character: HomeCharacterOutput{
					SupporterType: string(entity.SupporterTypeDog),
				},
				TodayWalk: &HomeTodayWalkOutput{
					WalkID: parsedUUID,
					Status: entity.WalkStatusWalking,
				},
			},
			wantErr: nil,
		},
		{
			name:   "まだ今日の散歩をしていない",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock, *mock.WalkRepositoryMock) {
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        parsedUUID,
								Nickname:  "test user",
								CreatedAt: now,
								UpdatedAt: now,
							}, nil
						},
					}, &mock.CharacterRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
							return &entity.Character{
								ID:            parsedUUID,
								UserID:        parsedUUID,
								SupporterType: entity.SupporterTypeDog,
								CreatedAt:     now,
								UpdatedAt:     now,
							}, nil
						},
					}, &mock.WalkRepositoryMock{
						FindTodayWalkByUserIDFunc: func(ctx context.Context, userID uuid.UUID, todayStart, tomorrowStart time.Time) (*entity.Walk, error) {
							return nil, nil
						},
					}
			},
			want: HomeOutput{
				Walker: HomeWalkerOutput{
					Nickname: "test user",
				},
				Character: HomeCharacterOutput{
					SupporterType: string(entity.SupporterTypeDog),
				},
				TodayWalk: nil,
			},
			wantErr: nil,
		},
		{
			name:   "Home一覧の取得を失敗する",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock, *mock.WalkRepositoryMock) {
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return nil, findErr
						},
					}, &mock.CharacterRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
							return nil, findErr
						},
					}, &mock.WalkRepositoryMock{
						FindTodayWalkByUserIDFunc: func(ctx context.Context, userID uuid.UUID, todayStart, tomorrowStart time.Time) (*entity.Walk, error) {
							return nil, findErr
						},
					}
			},
			want:    HomeOutput{},
			wantErr: findErr,
		},
		{
			name:   "Userの取得を失敗する",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock, *mock.WalkRepositoryMock) {
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return nil, findErr
						},
					}, &mock.CharacterRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
							return &entity.Character{SupporterType: entity.SupporterTypeDog}, nil
						},
					}, &mock.WalkRepositoryMock{
						FindTodayWalkByUserIDFunc: func(ctx context.Context, userID uuid.UUID, todayStart, tomorrowStart time.Time) (*entity.Walk, error) {
							return nil, nil
						},
					}
			},
			want:    HomeOutput{},
			wantErr: findErr,
		},
		{
			name:   "退会済ユーザーの場合ErrUnauthorizedを返す",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock, *mock.WalkRepositoryMock) {
				deletedAt := now
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        parsedUUID,
								Nickname:  "deleted user",
								CreatedAt: now,
								UpdatedAt: now,
								DeletedAt: &deletedAt,
							}, nil
						},
					}, &mock.CharacterRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
							return &entity.Character{SupporterType: entity.SupporterTypeDog}, nil
						},
					}, &mock.WalkRepositoryMock{
						FindTodayWalkByUserIDFunc: func(ctx context.Context, userID uuid.UUID, todayStart, tomorrowStart time.Time) (*entity.Walk, error) {
							return nil, nil
						},
					}
			},
			want:    HomeOutput{},
			wantErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, characterRepo, walkRepo := tt.mockFunc()
			uc := NewHomeUsecase(userRepo, characterRepo, walkRepo)
			got, err := uc.GetHome(ctx, tt.UserID)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
