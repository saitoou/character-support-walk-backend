package repository

import (
	"context"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	domainRepo "github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ domainRepo.UserRepository = (*UserRepositoryMock)(nil)

type UserRepositoryMock struct {
	// Read
	FindByIDFunc func(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	// Write
	CreateFunc          func(ctx context.Context, user entity.User) error
	UpdateFunc          func(ctx context.Context, userID uuid.UUID, nickname string, updatedAt time.Time) (*entity.User, error)
	UpdateDeletedAtFunc func(ctx context.Context, userID uuid.UUID, nickname string, deletedAt time.Time, updatedAt time.Time) error
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	if m.FindByIDFunc == nil {
		return nil, nil
	}
	return m.FindByIDFunc(ctx, userID)
}

func (m *UserRepositoryMock) Create(ctx context.Context, user entity.User) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, user)
}

func (m *UserRepositoryMock) Update(ctx context.Context, userID uuid.UUID, nickname string, updatedAt time.Time) (*entity.User, error) {
	if m.UpdateFunc == nil {
		return nil, nil
	}
	return m.UpdateFunc(ctx, userID, nickname, updatedAt)
}

func (m *UserRepositoryMock) UpdateDeletedAt(ctx context.Context, userID uuid.UUID, nickname string, deletedAt time.Time, updatedAt time.Time) error {
	if m.UpdateDeletedAtFunc == nil {
		return nil
	}
	return m.UpdateDeletedAtFunc(ctx, userID, nickname, deletedAt, updatedAt)
}
