package usecase

import (
	"context"
	"fmt"

	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
)

type WalkOptionOutput struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Title    string `json:"title"`
}

type WalkOptionListOutput struct {
	WalkOptions []WalkOptionOutput `json:"walk_options"`
}

type WalkOptionUsecase struct {
	repo repository.WalkOptionRepository
}

func NewWalkOptionUsecase(repo repository.WalkOptionRepository) *WalkOptionUsecase {
	return &WalkOptionUsecase{repo: repo}
}

func (uc WalkOptionUsecase) List(ctx context.Context) (WalkOptionListOutput, error) {

	res, err := uc.repo.SelectOption(ctx)
	if err != nil {
		return WalkOptionListOutput{}, fmt.Errorf("failed to find walk option list: %w", err)
	}

	outputs := make([]WalkOptionOutput, 0, len(res))

	for _, walkOption := range res {
		outputs = append(outputs, WalkOptionOutput{
			ID:       walkOption.ID,
			Category: walkOption.Category,
			Title:    walkOption.Title,
		})
	}

	return WalkOptionListOutput{
		WalkOptions: outputs,
	}, nil
}
