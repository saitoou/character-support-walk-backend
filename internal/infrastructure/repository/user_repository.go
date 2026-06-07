package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/apperr"
	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type UserRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"UserRepository.FindByID",
	)
	defer span.End()

	// query := `
	//   SELECT
	// 	id,
	// 	nickname,
	// 	created_at,
	// 	updated_at
	//   FROM users
	//   WHERE id = $1
	//     AND deleted_at IS NULL;
	// `
	// var user entity.User

	// db := r.executor(ctx)
	// row := db.QueryRowContext(ctx, query, userID)

	// if err := row.Scan(
	// 	&user.ID,
	// 	&user.Nickname,
	// 	&user.CreatedAt,
	// 	&user.UpdatedAt,
	// ); err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, nil
	// 	}
	// 	return nil, apperr.Wrap("failed to find user", err)
	// }

	// return &user, nil

	return nil, apperr.Wrap(
		"test repository error",
		fmt.Errorf("test error"),
	)
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	query := `
	  INSERT INTO users (
	    id,
		nickname,
		created_at,
		updated_at 
	  )
	  VALUES($1, $2, $3, $4)
	`
	db := r.executor(ctx)
	_, err := db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Nickname,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, userID uuid.UUID, nickname string, updatedAt time.Time) (*entity.User, error) {

	query := `
	  UPDATE users
	  SET
	    nickname = $1,
		updated_at = $2
	  WHERE id = $3
	    AND deleted_at IS NULL
	  RETURNING id, nickname, created_at, updated_at;
	`
	var user entity.User
	db := r.executor(ctx)
	err := db.QueryRowContext(ctx, query, nickname, updatedAt, userID).Scan(&user.ID, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update :%w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateDeletedAt(ctx context.Context, userID uuid.UUID, nickname string, deletedAt time.Time, updatedAt time.Time) error {
	query := `
	  UPDATE users
	  SET
	    nickname = $1,
		updated_at = $2,
		deleted_at = $3
	  WHERE id = $4
	    AND deleted_at IS NULL;
	`
	db := r.executor(ctx)

	result, err := db.ExecContext(ctx, query, nickname, updatedAt, deletedAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update user deleted_at :%w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) executor(ctx context.Context) database.DBTX {
	return database.GetExecutor(ctx, r.db)
}
