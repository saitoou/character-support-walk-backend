package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Nickname  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
