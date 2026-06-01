package entity

import (
	"time"

	"github.com/google/uuid"
)

type WalkStatus string

const (
	WalkStatusWalking   WalkStatus = "walking"
	WalkStatusCompleted WalkStatus = "completed"
	WalkStatusCanceled  WalkStatus = "canceled"
)

type Walk struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	WalkOptionID int
	Status       WalkStatus
	StartedAt    time.Time
	FinishedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
