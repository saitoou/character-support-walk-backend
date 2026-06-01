package repository

import (
	"context"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type CharacterRepository interface {
	// Read
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Character, error)
	// Write
	Create(ctx context.Context, character entity.Character) error
	UpdateByUserID(ctx context.Context, userID uuid.UUID, supporterType entity.SupporterType, updatedAt time.Time) (*entity.Character, error)
}
