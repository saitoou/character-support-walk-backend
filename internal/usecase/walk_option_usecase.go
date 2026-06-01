package usecase

import (
	"context"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
)

type WalkOptionUsecase struct {
	repo repository.WalkOptionRepository
}

func NewWalkOptionUsecase(repo repository.WalkOptionRepository) *WalkOptionUsecase {
	return &WalkOptionUsecase{repo: repo}
}

func (uc WalkOptionUsecase) List(ctx context.Context) ([]entity.WalkOption, error) {
	return uc.repo.SelectOption(ctx)
}
