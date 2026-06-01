package repository

import (
	"context"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	domainRepo "github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ domainRepo.WalkRepository = (*WalkRepositoryMock)(nil)

type WalkRepositoryMock struct {
	// Write
	CreateFunc       func(ctx context.Context, walk entity.Walk) error
	UpdateStatusFunc func(ctx context.Context, walk entity.Walk) error
	// Read
	FindByUserIDFunc          func(ctx context.Context, userID uuid.UUID) ([]entity.Walk, error)
	FindByUserIDAndIDFunc     func(ctx context.Context, userID uuid.UUID, walkID uuid.UUID) (*entity.Walk, error)
	FindTodayWalkByUserIDFunc func(ctx context.Context, userID uuid.UUID, todayStart time.Time, tomorrowStart time.Time) (*entity.Walk, error)
}

func (m *WalkRepositoryMock) Create(ctx context.Context, walk entity.Walk) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, walk)
}

func (m *WalkRepositoryMock) UpdateStatus(ctx context.Context, walk entity.Walk) error {
	if m.UpdateStatusFunc == nil {
		return nil
	}
	return m.UpdateStatusFunc(ctx, walk)
}

func (m *WalkRepositoryMock) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Walk, error) {
	if m.FindByUserIDFunc == nil {
		return nil, nil
	}
	return m.FindByUserIDFunc(ctx, userID)
}

func (m *WalkRepositoryMock) FindByUserIDAndID(ctx context.Context, userID uuid.UUID, walkID uuid.UUID) (*entity.Walk, error) {
	if m.FindByUserIDAndIDFunc == nil {
		return nil, nil
	}
	return m.FindByUserIDAndIDFunc(ctx, userID, walkID)
}

func (m *WalkRepositoryMock) FindTodayWalkByUserID(ctx context.Context, userID uuid.UUID, todayStart time.Time, tomorrowStart time.Time) (*entity.Walk, error) {
	if m.FindTodayWalkByUserIDFunc == nil {
		return nil, nil
	}
	return m.FindTodayWalkByUserIDFunc(ctx, userID, todayStart, tomorrowStart)
}
