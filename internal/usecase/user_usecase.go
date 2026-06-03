package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/google/uuid"
)

type WalkerOutput struct {
	UserID        uuid.UUID `json:"id"`
	Nickname      string    `json:"nickname"`
	SupporterType string    `json:"supporter_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserUsecase struct {
	userRepo      repository.UserRepository
	characterRepo repository.CharacterRepository
	authRepo      repository.AuthRepository
	tx            repository.Transaction
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	characterRepo repository.CharacterRepository,
	authRepo repository.AuthRepository,
	tx repository.Transaction,
) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, characterRepo: characterRepo, authRepo: authRepo, tx: tx}
}

func (uc *UserUsecase) FindByID(ctx context.Context, userID string) (WalkerOutput, error) {

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return WalkerOutput{}, fmt.Errorf("failed to parse uuid :%w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, parsedUserID)
	if err != nil {
		return WalkerOutput{}, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil || user.DeletedAt != nil {
		return WalkerOutput{}, fmt.Errorf("inactive user: %w", ErrUnauthorized)
	}

	character, err := uc.characterRepo.FindByUserID(ctx, parsedUserID)
	if err != nil {
		return WalkerOutput{}, fmt.Errorf("failed to find character: %w", err)
	}

	userProfile := WalkerOutput{
		UserID:        parsedUserID,
		Nickname:      user.Nickname,
		SupporterType: string(character.SupporterType),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	return userProfile, nil

}

func (uc *UserUsecase) UpdateByID(
	ctx context.Context,
	userID string,
	nickname string,
	supporterType string,
) (WalkerOutput, error) {

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return WalkerOutput{}, fmt.Errorf("failed to parse uuid :%v", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return WalkerOutput{}, fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	parsedSupporterType, err := entity.NewSupporterType(supporterType)
	if err != nil {
		return WalkerOutput{}, fmt.Errorf("invalid suppoter type: %w", err)
	}
	now := time.Now().UTC()

	var (
		updatedUser      *entity.User
		updatedCharacter *entity.Character
	)

	err = uc.tx.RunInTx(ctx, func(ctx context.Context) error {
		updatedUser, err = uc.userRepo.Update(ctx, parsedUserID, nickname, now)
		if err != nil {
			return fmt.Errorf("failed to updated user nickname: %w", err)
		}
		updatedCharacter, err = uc.characterRepo.UpdateByUserID(ctx, parsedUserID, parsedSupporterType, now)
		if err != nil {
			return fmt.Errorf("failed to updated character supporter_type: %w", err)
		}
		return nil
	})
	if err != nil {
		return WalkerOutput{}, fmt.Errorf("failed to update user profile: %w", err)
	}

	return WalkerOutput{
		UserID:        updatedUser.ID,
		Nickname:      updatedUser.Nickname,
		SupporterType: string(updatedCharacter.SupporterType),
		CreatedAt:     updatedUser.CreatedAt,
		UpdatedAt:     updatedUser.UpdatedAt,
	}, nil
}

func (uc *UserUsecase) Deactivate(ctx context.Context, userID string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("failed to parse uuid: %w", err)
	}
	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}
	now := time.Now().UTC()
	err = uc.tx.RunInTx(ctx, func(ctx context.Context) error {
		err := uc.userRepo.UpdateDeletedAt(ctx, parsedUserID, fmt.Sprintf("deleted_user_%s", parsedUserID.String()[:8]), now, now)
		if err != nil {
			return fmt.Errorf("failed to update deleted_at: %w", err)
		}

		err = uc.authRepo.DeleteAuthIdentitiesByUserID(ctx, parsedUserID)
		if err != nil {
			return fmt.Errorf("failed to deleted auth identities: %w", err)
		}

		err = uc.authRepo.UpdateRefreshTokenRevokedAtByUserID(ctx, parsedUserID, now)
		if err != nil {
			return fmt.Errorf("failed to revoke refresh tokens: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}
