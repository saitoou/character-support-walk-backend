package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SupporterType string

const (
	SupporterTypeDog SupporterType = "dog"
	SupporterTypeCat SupporterType = "cat"
)

type Character struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	SupporterType SupporterType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSupporterType(v string) (SupporterType, error) {
	switch v {
	case string(SupporterTypeDog):
		return SupporterTypeDog, nil
	case string(SupporterTypeCat):
		return SupporterTypeCat, nil
	default:
		return "", fmt.Errorf("invalid supporter type: %s", v)
	}
}
