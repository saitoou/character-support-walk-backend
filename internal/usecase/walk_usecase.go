package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
)

type StartWalkOutput struct {
	WalkID uuid.UUID
	Status string
}

type WalkOutput struct {
	ID           uuid.UUID  `json:"id"`
	WalkOptionID int        `json:"walk_option_id"`
	Status       string     `json:"status"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type WalkListOutput struct {
	Walks []WalkOutput `json:"walks"`
}

type WalkUsecase struct {
	walkRepo repository.WalkRepository
	userRepo repository.UserRepository
}

func NewWalkUsecase(walkRepo repository.WalkRepository, userRepo repository.UserRepository) *WalkUsecase {
	return &WalkUsecase{walkRepo: walkRepo, userRepo: userRepo}
}

func (uc *WalkUsecase) Start(ctx context.Context, walkOptionId int, userID string) (StartWalkOutput, error) {

	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"WalkUsecase.Start",
	)
	defer span.End()

	now := time.Now().UTC()
	walkUUID, err := uuid.NewV7()
	if err != nil {
		return StartWalkOutput{}, fmt.Errorf("failed to generate uuid: %v", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return StartWalkOutput{}, fmt.Errorf("failed to parse uuid :%w", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return StartWalkOutput{}, fmt.Errorf("ensure active user: %w", err)
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
		return StartWalkOutput{}, fmt.Errorf("failed to start walk: %w", err)
	}
	return StartWalkOutput{
		WalkID: walkUUID,
		Status: string(walk.Status),
	}, nil
}

func (uc *WalkUsecase) GetWalks(ctx context.Context, userID string) (WalkListOutput, error) {

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return WalkListOutput{}, fmt.Errorf("failed to parse uuid :%v", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return WalkListOutput{}, fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	result, err := uc.walkRepo.FindByUserID(ctx, parsedUserID)
	if err != nil {
		return WalkListOutput{}, fmt.Errorf("failed to get walks: %w", err)
	}

	outputs := make([]WalkOutput, 0, len(result))

	for _, walk := range result {
		outputs = append(outputs, WalkOutput{
			ID:           walk.ID,
			WalkOptionID: walk.WalkOptionID,
			Status:       string(walk.Status),
			StartedAt:    walk.StartedAt,
			FinishedAt:   walk.FinishedAt,
			CreatedAt:    walk.CreatedAt,
			UpdatedAt:    walk.UpdatedAt,
		})
	}

	return WalkListOutput{Walks: outputs}, nil
}

func (uc *WalkUsecase) GetWalking(ctx context.Context, userID string, walkID string) (WalkOutput, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return WalkOutput{}, fmt.Errorf("failed to parse uuid :%v", err)
	}

	if err := EnsureActiveUser(ctx, uc.userRepo, parsedUserID); err != nil {
		return WalkOutput{}, fmt.Errorf("ensure active user: %w", ErrUnauthorized)
	}

	parsedWalkID, err := uuid.Parse(walkID)
	if err != nil {
		return WalkOutput{}, fmt.Errorf("failed to parse uuid :%v", err)
	}

	result, err := uc.walkRepo.FindByUserIDAndID(ctx, parsedUserID, parsedWalkID)
	if err != nil {
		return WalkOutput{}, fmt.Errorf("failed to get walks: %w", err)
	}

	output := WalkOutput{
		ID:           result.ID,
		WalkOptionID: result.WalkOptionID,
		Status:       string(result.Status),
		StartedAt:    result.StartedAt,
		FinishedAt:   result.FinishedAt,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return output, nil
}

func (uc *WalkUsecase) UpdateComplete(ctx context.Context, userID, walkID string) error {
	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"WalkUsecase.UpdateComplete",
	)
	defer span.End()

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
	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"WalkUsecase.UpdateCancel",
	)
	defer span.End()

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
