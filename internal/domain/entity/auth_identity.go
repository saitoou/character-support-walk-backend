package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	AuthProviderGoogle AuthProvider = "google"
)

type AuthIdentity struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       AuthProvider
	ProviderUserID string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
