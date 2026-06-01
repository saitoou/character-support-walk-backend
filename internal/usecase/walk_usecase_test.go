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

func Test_WalkUsecase_Start(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	// parsedUUID := uuid.MustParse(id)
	insertErr := errors.New("failed to start walk")
	// now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name         string
		WalkOptionID int
		UserID       string
		mockFunc     func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock)
		wantStatus   string
		wantErr      error
	}{
		{
			name:         "Start Walkが成功する",
			WalkOptionID: 1,
			UserID:       id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						CreateFunc: func(ctx context.Context, walk entity.Walk) error {
							return nil
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantStatus: string(entity.WalkStatusWalking),
			wantErr:    nil,
		},
		{
			name:         "Start Walkが失敗する",
			WalkOptionID: 1,
			UserID:       id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						CreateFunc: func(ctx context.Context, walk entity.Walk) error {
							return insertErr
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantStatus: "",
			wantErr:    insertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walkRepo, userRepo := tt.mockFunc()
			uc := NewWalkUsecase(walkRepo, userRepo)
			got, err := uc.Start(ctx, tt.WalkOptionID, tt.UserID)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func Test_WalkUsecase_GetWalks(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	findErr := errors.New("failed to get walks")
	now := time.Now().UTC().Truncate(time.Millisecond)

	walks := []entity.Walk{
		{
			ID:           parsedUUID,
			UserID:       parsedUUID,
			WalkOptionID: 1,
			Status:       entity.WalkStatusWalking,
			StartedAt:    now,
			FinishedAt:   nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           parsedUUID,
			UserID:       parsedUUID,
			WalkOptionID: 1,
			Status:       entity.WalkStatusWalking,
			StartedAt:    now,
			FinishedAt:   nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name     string
		want     []entity.Walk
		UserID   string
		mockFunc func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock)
		wantErr  error
	}{
		{
			name:   "Walk一覧取得",
			want:   walks,
			UserID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) ([]entity.Walk, error) {
							return walks, nil
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: nil,
		},
		{
			name:   "Walk一覧取得失敗",
			want:   nil,
			UserID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) ([]entity.Walk, error) {
							return nil, findErr
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walkRepo, userRepo := tt.mockFunc()
			uc := NewWalkUsecase(walkRepo, userRepo)
			got, err := uc.GetWalks(ctx, tt.UserID)
			assert.Equal(t, got, tt.want)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkUsecase_GetWalking(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	findErr := errors.New("failed to get walks")
	now := time.Now().UTC().Truncate(time.Millisecond)

	walk := &entity.Walk{
		ID:           parsedUUID,
		UserID:       parsedUUID,
		WalkOptionID: 1,
		Status:       entity.WalkStatusWalking,
		StartedAt:    now,
		FinishedAt:   nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name     string
		want     *entity.Walk
		userID   string
		walkID   string
		mockFunc func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock)
		wantErr  error
	}{
		{
			name:   "Walk一覧取得",
			want:   walk,
			userID: id,
			walkID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDAndIDFunc: func(ctx context.Context, userID, walkID uuid.UUID) (*entity.Walk, error) {
							return walk, nil
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: nil,
		},
		{
			name:   "Walk一覧取得失敗",
			want:   nil,
			userID: id,
			walkID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDAndIDFunc: func(ctx context.Context, userID, walkID uuid.UUID) (*entity.Walk, error) {
							return nil, findErr
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: findErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walkRepo, userRepo := tt.mockFunc()
			uc := NewWalkUsecase(walkRepo, userRepo)
			got, err := uc.GetWalking(ctx, tt.userID, tt.walkID)
			assert.Equal(t, got, tt.want)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkUsecase_UpdateComplete(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	updateErr := errors.New("failed to update walk")
	now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name     string
		userID   string
		walkID   string
		mockFunc func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock)
		wantErr  error
	}{
		{
			name:   "WalkのStatus更新成功",
			userID: id,
			walkID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDAndIDFunc: func(ctx context.Context, userID, walkID uuid.UUID) (*entity.Walk, error) {
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
						UpdateStatusFunc: func(ctx context.Context, updatedWalk entity.Walk) error {
							return nil
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: nil,
		},
		{
			name:   "WalkのStatus更新の失敗",
			userID: id,
			walkID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDAndIDFunc: func(ctx context.Context, userID, walkID uuid.UUID) (*entity.Walk, error) {
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
						UpdateStatusFunc: func(ctx context.Context, updatedWalk entity.Walk) error {
							return updateErr
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: updateErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walkRepo, userRepo := tt.mockFunc()
			uc := NewWalkUsecase(walkRepo, userRepo)
			err := uc.UpdateComplete(ctx, tt.userID, tt.walkID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_WalkUsecase_UpdateCancel(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	updateErr := errors.New("failed to update walk")
	now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name     string
		userID   string
		walkID   string
		mockFunc func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock)
		wantErr  error
	}{
		{
			name:   "WalkのStatus更新成功",
			userID: id,
			walkID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDAndIDFunc: func(ctx context.Context, userID, walkID uuid.UUID) (*entity.Walk, error) {
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
						UpdateStatusFunc: func(ctx context.Context, updatedWalk entity.Walk) error {
							return nil
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: nil,
		},
		{
			name:   "WalkのStatus更新の失敗",
			userID: id,
			walkID: id,
			mockFunc: func() (*mock.WalkRepositoryMock, *mock.UserRepositoryMock) {
				return &mock.WalkRepositoryMock{
						FindByUserIDAndIDFunc: func(ctx context.Context, userID, walkID uuid.UUID) (*entity.Walk, error) {
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
						UpdateStatusFunc: func(ctx context.Context, updatedWalk entity.Walk) error {
							return updateErr
						},
					},
					&mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        userID,
								Nickname:  "test user",
								DeletedAt: nil,
							}, nil
						},
					}
			},
			wantErr: updateErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walkRepo, userRepo := tt.mockFunc()
			uc := NewWalkUsecase(walkRepo, userRepo)
			err := uc.UpdateCancel(ctx, tt.userID, tt.walkID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
