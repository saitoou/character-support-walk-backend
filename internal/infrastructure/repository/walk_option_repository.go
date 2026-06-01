package repository

import (
	"context"
	"fmt"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/database"
)

type WalkOptionRepository struct {
	db database.DBTX
}

func NewWalkOptionRepository(db database.DBTX) *WalkOptionRepository {
	return &WalkOptionRepository{db: db}
}

func (r *WalkOptionRepository) SelectOption(ctx context.Context) ([]entity.WalkOption, error) {

	db := r.executor(ctx)
	rows, err := db.QueryContext(ctx, "SELECT id, category, title FROM walk_options ORDER BY id ASC;")
	if err != nil {
		return nil, fmt.Errorf("failed to select data: %w", err)
	}

	defer rows.Close()
	var options []entity.WalkOption

	for rows.Next() {
		var option entity.WalkOption

		if err := rows.Scan(&option.ID, &option.Category, &option.Title); err != nil {
			return nil, fmt.Errorf("failed to apply work options: %v", err)
		}
		fmt.Printf("ID: %s, Category: %s, Title: %s", option.ID, option.Category, option.Title)

		options = append(options, option)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return options, nil
}

func (r *WalkOptionRepository) executor(ctx context.Context) database.DBTX {
	return database.GetExecutor(ctx, r.db)
}
