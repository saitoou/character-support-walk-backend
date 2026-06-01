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
)

type WalkRepository struct {
	db database.DBTX
}

func NewWalkRepository(db database.DBTX) *WalkRepository {
	return &WalkRepository{db: db}
}

func (r *WalkRepository) Create(ctx context.Context, walk entity.Walk) error {
	query := `
	  INSERT INTO walks (
	    id,
		user_id,
		walk_option_id,
		status,
		started_at,
		finished_at,
		created_at,
		updated_at 
	  )
	  VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`

	db := r.executor(ctx)
	_, err := db.ExecContext(
		ctx,
		query,
		walk.ID,
		walk.UserID,
		walk.WalkOptionID,
		walk.Status,
		walk.StartedAt,
		walk.FinishedAt,
		walk.CreatedAt,
		walk.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil

}

func (r *WalkRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Walk, error) {

	query := `
	  SElECT
	    id,
		user_id,
		walk_option_id,
		status,
		started_at,
		finished_at,
		created_at,
		updated_at
	  FROM walks
	  WHERE user_id = $1
	  ORDER BY created_at DESC
	`

	db := r.executor(ctx)
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find walks: %w", err)
	}

	defer rows.Close()
	var walks []entity.Walk

	for rows.Next() {
		var walk entity.Walk
		if err := rows.Scan(
			&walk.ID,
			&walk.UserID,
			&walk.WalkOptionID,
			&walk.Status,
			&walk.StartedAt,
			&walk.FinishedAt,
			&walk.CreatedAt,
			&walk.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to apply walk: %w", err)
		}

		walks = append(walks, walk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return walks, nil
}

func (r *WalkRepository) FindByUserIDAndID(ctx context.Context, userID uuid.UUID, walkID uuid.UUID) (*entity.Walk, error) {

	query := `
	  SELECT
	    id,
		user_id,
		walk_option_id,
		status,
		started_at,
		finished_at,
		created_at,
		updated_at
	  FROM walks
	  WHERE user_id = $1
	    AND ID = $2
	`

	var walk entity.Walk

	db := r.executor(ctx)
	row := db.QueryRowContext(ctx, query, userID, walkID)

	if err := row.Scan(
		&walk.ID,
		&walk.UserID,
		&walk.WalkOptionID,
		&walk.Status,
		&walk.StartedAt,
		&walk.FinishedAt,
		&walk.CreatedAt,
		&walk.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to walk by user id and id error: %w", err)
	}

	return &walk, nil
}

func (r *WalkRepository) UpdateStatus(ctx context.Context, walk entity.Walk) error {
	query := `
	  UPDATE walks
	  SET
	    status = $1,
		finished_at = $2,
		updated_at = $3
	  WHERE id = $4
	    AND user_id = $5
		AND status = 'walking'
	`

	result, err := r.db.ExecContext(ctx, query, walk.Status, walk.FinishedAt, walk.UpdatedAt, walk.ID, walk.UserID)
	if err != nil {
		return fmt.Errorf("failed to update :%w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("walk not found or not owned by user")
	}

	return nil

}

func (r *WalkRepository) FindTodayWalkByUserID(ctx context.Context, userID uuid.UUID, todayStart time.Time, tomorrowStart time.Time) (*entity.Walk, error) {

	query := `
	  SElECT
	    id,
		user_id,
		walk_option_id,
		status,
		started_at,
		finished_at,
		created_at,
		updated_at
	  FROM walks
	  WHERE user_id = $1
	    AND created_at >= $2
		AND created_at < $3
	  ORDER BY
	    CASE
		  WHEN status = 'completed' THEN 1
		  WHEN status = 'walking' THEN 2
		  WHEN status = 'canceled' THEN 3
		  ELSE 4
		END ASC,
		created_at DESC
	  LIMIT 1
	`
	var walk entity.Walk

	db := r.executor(ctx)
	row := db.QueryRowContext(ctx, query, userID, todayStart, tomorrowStart)

	if err := row.Scan(
		&walk.ID,
		&walk.UserID,
		&walk.WalkOptionID,
		&walk.Status,
		&walk.StartedAt,
		&walk.FinishedAt,
		&walk.CreatedAt,
		&walk.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to walk by user id and id error: %w", err)
	}

	return &walk, nil
}

func (r *WalkRepository) executor(ctx context.Context) database.DBTX {
	return database.GetExecutor(ctx, r.db)
}
