package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
)

type StartWalkOutput struct {
	WalkID uuid.UUID
	Status string
}

type WalkUsecase struct {
	walkRepo repository.WalkRepository
	userRepo repository.UserRepository
}

func NewWalkUsecase(walkRepo repository.WalkRepository, userRepo repository.UserRepository) *WalkUsecase {
	return &WalkUsecase{walkRepo: walkRepo, userRepo: userRepo}
}

func (uc *WalkUsecase) Start(ctx context.Context, walkOptionId int, userID string) (*StartWalkOutput, error) {

	now := time.Now().UTC()
	walkUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %v", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid :%w", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return nil, fmt.Errorf("ensure active user: %w", err)
	}

	walk := entity.Walk{
		ID:           walkUUID,
		UserID:       parsedUserID,
		WalkOptionID: walkOptionId,
		Status:       entity.WalkStatusWalking,
		StartedAt:    now,
		FinishedAt:   nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.walkRepo.Create(ctx, walk); err != nil {
		return nil, fmt.Errorf("failed to start walk: %w", err)
	}
	return &StartWalkOutput{
		WalkID: walkUUID,
		Status: string(walk.Status),
	}, nil
}

func (uc *WalkUsecase) GetWalks(ctx context.Context, userID string) ([]entity.Walk, error) {

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid :%v", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return nil, fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	result, err := uc.walkRepo.FindByUserID(ctx, parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get walks: %w", err)
	}

	return result, nil
}

func (uc *WalkUsecase) GetWalking(ctx context.Context, userID string, walkID string) (*entity.Walk, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid :%v", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return nil, fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	parsedWalkID, err := uuid.Parse(walkID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid :%v", err)
	}

	result, err := uc.walkRepo.FindByUserIDAndID(ctx, parsedUserID, parsedWalkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get walks: %w", err)
	}

	return result, nil
}

func (uc *WalkUsecase) UpdateComplete(ctx context.Context, userID, walkID string) error {

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("failed to parse uuid :%w", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	parsedWalkID, err := uuid.Parse(walkID)
	if err != nil {
		return fmt.Errorf("failed to parse uuid :%w", err)
	}

	result, err := uc.walkRepo.FindByUserIDAndID(ctx, parsedUserID, parsedWalkID)
	if err != nil {
		return fmt.Errorf("faileed to find walks by userID and ID: %w", err)
	}

	if result.Status != entity.WalkStatusWalking {
		return fmt.Errorf("walk status is not walking")
	}

	now := time.Now().UTC()

	result.Status = entity.WalkStatusCompleted
	result.FinishedAt = &now
	result.UpdatedAt = now

	err = uc.walkRepo.UpdateStatus(ctx, *result)
	if err != nil {
		return fmt.Errorf("failed to update completed walk :%w", err)
	}

	return nil
}

func (uc *WalkUsecase) UpdateCancel(ctx context.Context, userID, walkID string) error {

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("failed to parse uuid :%w", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	parsedWalkID, err := uuid.Parse(walkID)
	if err != nil {
		return fmt.Errorf("failed to parse uuid :%w", err)
	}

	result, err := uc.walkRepo.FindByUserIDAndID(ctx, parsedUserID, parsedWalkID)
	if err != nil {
		return fmt.Errorf("faileed to find walks by userID and ID: %w", err)
	}

	if result.Status != entity.WalkStatusWalking {
		return fmt.Errorf("walk status is not walking")
	}

	now := time.Now().UTC()

	result.Status = entity.WalkStatusCanceled
	result.FinishedAt = nil
	result.UpdatedAt = now

	err = uc.walkRepo.UpdateStatus(ctx, *result)
	if err != nil {
		return fmt.Errorf("failed to update canceled walk :%w", err)
	}

	return nil
}
