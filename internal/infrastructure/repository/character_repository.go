package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type CharacterRepository struct {
	db database.DBTX
}

func NewCharacterRepository(db database.DBTX) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (r *CharacterRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Character, error) {

	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"CharacterRepository.FindByID",
	)
	defer span.End()
	query := `
	  SElECT
		id,
		user_id,
		supporter_type,
		created_at,
		updated_at
	  FROM characters
	  WHERE user_id = $1
	`

	var character entity.Character

	db := r.executor(ctx)
	row := db.QueryRowContext(ctx, query, userID)

	if err := row.Scan(
		&character.ID,
		&character.UserID,
		&character.SupporterType,
		&character.CreatedAt,
		&character.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	return &character, nil
}

func (r *CharacterRepository) UpdateByUserID(ctx context.Context, userID uuid.UUID, supporterType entity.SupporterType, updatedAt time.Time) (*entity.Character, error) {

	query := `
	  UPDATE characters
	  SET
	    supporter_type = $1,
		updated_at = $2
	  WHERE user_id = $3
	  RETURNING user_id, supporter_type, updated_at;
	`
	var character entity.Character
	db := r.executor(ctx)
	err := db.QueryRowContext(ctx, query, supporterType, updatedAt, userID).Scan(&character.UserID, &character.SupporterType, &character.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to update :%w", err)
	}
	return &character, nil
}

func (r *CharacterRepository) Create(ctx context.Context, character entity.Character) error {
	query := `
	  INSERT INTO characters (
	    id,
		user_id,
		supporter_type,
		created_at,
		updated_at 
	  )
	  VALUES($1, $2, $3, $4, $5)
	`

	db := r.executor(ctx)
	_, err := db.ExecContext(
		ctx,
		query,
		character.ID,
		character.UserID,
		character.SupporterType,
		character.CreatedAt,
		character.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert character: %w", err)
	}

	return nil
}

func (r *CharacterRepository) executor(ctx context.Context) database.DBTX {
	return database.GetExecutor(ctx, r.db)
}
