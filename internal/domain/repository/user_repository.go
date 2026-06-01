package repository

import (
	"context"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	// Read
	FindByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	// Write
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, userID uuid.UUID, nickname string, updatedAt time.Time) (*entity.User, error)
	UpdateDeletedAt(ctx context.Context, userID uuid.UUID, nickname string, deletedAt time.Time, updatedAt time.Time) error
}
