package usecase

import (
	"context"
	"fmt"

	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

func EnsureActiveUser(ctx context.Context, userRepo repository.UserRepository, userID uuid.UUID) error {

	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"Usecase.EnsureActiveUser",
	)
	defer span.End()

	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil || user.DeletedAt != nil {
		return fmt.Errorf("inactive user: %w", ErrUnauthorized)
	}

	return nil
}
