package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	mock "github.com/chocoko/character-support-walk-backend/internal/test/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UserUsecase_FindByID(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	// findErr := errors.New("failed to find walker")
	now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name     string
		UserID   string
		mockFunc func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock)
		want     WalkerOutput
		wantErr  error
	}{
		{
			name:   "Walkerを取得する",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock) {
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
					}
			},
			want: WalkerOutput{
				UserID:        parsedUUID,
				Nickname:      "test user",
				SupporterType: string(entity.SupporterTypeDog),
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			wantErr: nil,
		},
		{
			name:   "Walkerの取得に失敗する",
			UserID: id,
			mockFunc: func() (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock) {
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return nil, ErrUnauthorized
						},
					}, &mock.CharacterRepositoryMock{
						FindByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
							return nil, nil
						},
					}
			},
			want:    WalkerOutput{},
			wantErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tx := mock.NewMockTransaction(ctrl)
			authRepo := mock.NewMockAuthRepository(ctrl)
			userRepo, characterRepo := tt.mockFunc()
			uc := NewUserUsecase(userRepo, characterRepo, authRepo, tx)
			got, err := uc.FindByID(ctx, tt.UserID)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_UserUsecase_UpdateByID(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	updateErr := errors.New("failed to update walker")
	now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name          string
		UserID        string
		nickname      string
		supporterType string
		mockFunc      func(tx *mock.MockTransaction) (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock)
		want          WalkerOutput
		wantErr       error
	}{
		{
			name:          "Walkerを更新する",
			UserID:        id,
			nickname:      "test user",
			supporterType: string(entity.SupporterTypeDog),
			mockFunc: func(tx *mock.MockTransaction) (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        parsedUUID,
								Nickname:  "test user",
								CreatedAt: now,
								UpdatedAt: now,
								DeletedAt: nil,
							}, nil
						},
						UpdateFunc: func(ctx context.Context, userID uuid.UUID, nickname string, updatedAt time.Time) (*entity.User, error) {
							return &entity.User{
								ID:        parsedUUID,
								Nickname:  "test user",
								CreatedAt: now,
								UpdatedAt: now,
							}, nil
						},
					}, &mock.CharacterRepositoryMock{
						UpdateByUserIDFunc: func(ctx context.Context, userID uuid.UUID, supporterType entity.SupporterType, updatedAt time.Time) (*entity.Character, error) {
							return &entity.Character{
								ID:            parsedUUID,
								UserID:        parsedUUID,
								SupporterType: entity.SupporterTypeDog,
								CreatedAt:     now,
								UpdatedAt:     now,
							}, nil
						},
					}
			},
			want: WalkerOutput{
				UserID:        parsedUUID,
				Nickname:      "test user",
				SupporterType: string(entity.SupporterTypeDog),
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			wantErr: nil,
		},
		{
			name:          "Walkerの取得に失敗する",
			UserID:        id,
			nickname:      "test user",
			supporterType: string(entity.SupporterTypeDog),
			mockFunc: func(tx *mock.MockTransaction) (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				return &mock.UserRepositoryMock{
						FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
							return &entity.User{
								ID:        parsedUUID,
								Nickname:  "test user",
								CreatedAt: now,
								UpdatedAt: now,
								DeletedAt: nil,
							}, nil
						},
						UpdateFunc: func(ctx context.Context, userID uuid.UUID, nickname string, updatedAt time.Time) (*entity.User, error) {
							return nil, updateErr
						},
					}, &mock.CharacterRepositoryMock{
						UpdateByUserIDFunc: func(ctx context.Context, userID uuid.UUID, supporterType entity.SupporterType, updatedAt time.Time) (*entity.Character, error) {
							return nil, updateErr
						},
					}
			},
			want:    WalkerOutput{},
			wantErr: updateErr,
		},
		{
			name:          "退会済みユーザーの場合ErrUnauthorizedを返す",
			UserID:        id,
			nickname:      "deleted user",
			supporterType: string(entity.SupporterTypeDog),
			mockFunc: func(tx *mock.MockTransaction) (*mock.UserRepositoryMock, *mock.CharacterRepositoryMock) {
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
					},
					&mock.CharacterRepositoryMock{}
			},
			want:    WalkerOutput{},
			wantErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tx := mock.NewMockTransaction(ctrl)
			authRepo := mock.NewMockAuthRepository(ctrl)
			userRepo, characterRepo := tt.mockFunc(tx)
			uc := NewUserUsecase(userRepo, characterRepo, authRepo, tx)
			got, err := uc.UpdateByID(ctx, tt.UserID, tt.nickname, tt.supporterType)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_UserUsecase_Deactivate(t *testing.T) {
	ctx := context.Background()
	id := "019e1cd3-8194-7a36-816b-2f38206ca52d"
	parsedUUID := uuid.MustParse(id)
	updateDeleteAuthIdenetityErr := errors.New("failed to delete auth identity")
	upateRevorkErr := errors.New("failed to update revoke_at")
	updateDeletedAtErr := errors.New("failed to update deleted_at")
	now := time.Now().UTC().Truncate(time.Millisecond)

	tests := []struct {
		name     string
		UserID   string
		nickname string
		mockFunc func(tx *mock.MockTransaction, authRepo *mock.MockAuthRepository) (userRepo *mock.UserRepositoryMock)
		wantErr  error
	}{
		{
			name:     "退会したユーザーの削除が成功する",
			UserID:   id,
			nickname: fmt.Sprintf("deleted_user_%s", id[:8]),
			mockFunc: func(tx *mock.MockTransaction, authRepo *mock.MockAuthRepository) (userRepo *mock.UserRepositoryMock) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				authRepo.EXPECT().
					DeleteAuthIdentitiesByUserID(gomock.Any(), parsedUUID).
					Return(nil)

				authRepo.EXPECT().UpdateRefreshTokenRevokedAtByUserID(gomock.Any(), parsedUUID, gomock.Any()).
					Return(nil)

				return &mock.UserRepositoryMock{
					UpdateDeletedAtFunc: func(ctx context.Context, userID uuid.UUID, nickname string, deletedAt, updatedAt time.Time) error {
						return nil
					},
					FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
						return &entity.User{
							ID:        parsedUUID,
							Nickname:  fmt.Sprintf("deleted_user_%s", id[:8]),
							CreatedAt: now,
							UpdatedAt: now,
							DeletedAt: nil,
						}, nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name:     "退会する際、AuthIdentityの削除に失敗する",
			UserID:   id,
			nickname: fmt.Sprintf("deleted_user_%s", id[:8]),
			mockFunc: func(tx *mock.MockTransaction, authRepo *mock.MockAuthRepository) (userRepo *mock.UserRepositoryMock) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				authRepo.EXPECT().
					DeleteAuthIdentitiesByUserID(gomock.Any(), parsedUUID).
					Return(updateDeleteAuthIdenetityErr)

				return &mock.UserRepositoryMock{
					UpdateDeletedAtFunc: func(ctx context.Context, userID uuid.UUID, nickname string, deletedAt, updatedAt time.Time) error {
						return nil
					},
					FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
						return &entity.User{
							ID:        parsedUUID,
							Nickname:  fmt.Sprintf("deleted_user_%s", id[:8]),
							CreatedAt: now,
							UpdatedAt: now,
							DeletedAt: nil,
						}, nil
					},
				}
			},
			wantErr: updateDeleteAuthIdenetityErr,
		},
		{
			name:     "退会する際、RefreshTokenのRevokedAtの更新に失敗する",
			UserID:   id,
			nickname: fmt.Sprintf("deleted_user_%s", id[:8]),
			mockFunc: func(tx *mock.MockTransaction, authRepo *mock.MockAuthRepository) (userRepo *mock.UserRepositoryMock) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				authRepo.EXPECT().
					DeleteAuthIdentitiesByUserID(gomock.Any(), parsedUUID).
					Return(nil)

				authRepo.EXPECT().UpdateRefreshTokenRevokedAtByUserID(gomock.Any(), parsedUUID, gomock.Any()).
					Return(upateRevorkErr)

				return &mock.UserRepositoryMock{
					UpdateDeletedAtFunc: func(ctx context.Context, userID uuid.UUID, nickname string, deletedAt, updatedAt time.Time) error {
						return nil
					},
					FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
						return &entity.User{
							ID:        parsedUUID,
							Nickname:  fmt.Sprintf("deleted_user_%s", id[:8]),
							CreatedAt: now,
							UpdatedAt: now,
							DeletedAt: nil,
						}, nil
					},
				}
			},
			wantErr: upateRevorkErr,
		},
		{
			name:     "退会するメンバーのdeleted_atの更新に失敗する",
			UserID:   id,
			nickname: fmt.Sprintf("deleted_user_%s", id[:8]),
			mockFunc: func(tx *mock.MockTransaction, authRepo *mock.MockAuthRepository) (userRepo *mock.UserRepositoryMock) {
				tx.EXPECT().RunInTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				return &mock.UserRepositoryMock{
					FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
						return &entity.User{
							ID:        parsedUUID,
							Nickname:  fmt.Sprintf("deleted_user_%s", id[:8]),
							CreatedAt: now,
							UpdatedAt: now,
							DeletedAt: nil,
						}, nil
					},
					UpdateDeletedAtFunc: func(ctx context.Context, userID uuid.UUID, nickname string, deletedAt, updatedAt time.Time) error {
						return updateDeletedAtErr
					},
				}
			},
			wantErr: updateDeletedAtErr,
		},
		{
			name:     "退会済みユーザーの場合ErrUnauthorizedを返す",
			UserID:   id,
			nickname: fmt.Sprintf("deleted_user_%s", id[:8]),
			mockFunc: func(tx *mock.MockTransaction, authRepo *mock.MockAuthRepository) (userRepo *mock.UserRepositoryMock) {
				deletedAt := now
				return &mock.UserRepositoryMock{
					FindByIDFunc: func(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
						return &entity.User{
							ID:        parsedUUID,
							Nickname:  fmt.Sprintf("deleted_user_%s", id[:8]),
							CreatedAt: now,
							UpdatedAt: now,
							DeletedAt: &deletedAt,
						}, nil
					},
				}
			},
			wantErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tx := mock.NewMockTransaction(ctrl)
			authRepo := mock.NewMockAuthRepository(ctrl)
			characterRepo := &mock.CharacterRepositoryMock{}
			userRepo := tt.mockFunc(tx, authRepo)
			uc := NewUserUsecase(userRepo, characterRepo, authRepo, tx)
			err := uc.Deactivate(ctx, tt.UserID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
