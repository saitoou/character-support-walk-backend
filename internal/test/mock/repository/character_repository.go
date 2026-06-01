package repository

import (
	"context"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	domainRepo "github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ domainRepo.CharacterRepository = (*CharacterRepositoryMock)(nil)

type CharacterRepositoryMock struct {
	// Read
	FindByUserIDFunc func(ctx context.Context, userID uuid.UUID) (*entity.Character, error)
	// Write
	CreateFunc         func(ctx context.Context, character entity.Character) error
	UpdateByUserIDFunc func(ctx context.Context, userID uuid.UUID, supporterType entity.SupporterType, updatedAt time.Time) (*entity.Character, error)
}

func (m *CharacterRepositoryMock) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {
	if m.FindByUserIDFunc == nil {
		return nil, nil
	}
	return m.FindByUserIDFunc(ctx, userID)
}

func (m *CharacterRepositoryMock) Create(ctx context.Context, character entity.Character) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, character)
}

func (m *CharacterRepositoryMock) UpdateByUserID(ctx context.Context, userID uuid.UUID, supporterType entity.SupporterType, updatedAt time.Time) (*entity.Character, error) {
	if m.UpdateByUserIDFunc == nil {
		return nil, nil
	}
	return m.UpdateByUserIDFunc(ctx, userID, supporterType, updatedAt)
}
