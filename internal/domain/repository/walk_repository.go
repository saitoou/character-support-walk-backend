package repository

import (
	"context"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type WalkRepository interface {
	// Write
	Create(ctx context.Context, walk entity.Walk) error
	UpdateStatus(ctx context.Context, walk entity.Walk) error
	// Read
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Walk, error)
	FindByUserIDAndID(ctx context.Context, userID uuid.UUID, walkID uuid.UUID) (*entity.Walk, error)
	FindTodayWalkByUserID(ctx context.Context, userID uuid.UUID, todayStart time.Time, tomorrowStart time.Time) (*entity.Walk, error)
}
