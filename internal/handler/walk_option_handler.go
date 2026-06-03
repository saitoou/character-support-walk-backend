package handler

import (
	"context"
	"net/http"

	"github.com/chocoko/character-support-walk-backend/gen/openapi/v1"
	"github.com/chocoko/character-support-walk-backend/internal/usecase"
	"github.com/labstack/echo/v5"
)

type WalkOptionUsecase interface {
	List(ctx context.Context) (usecase.WalkOptionListOutput, error)
}

type WalkOptionHandler struct {
	uc *usecase.WalkOptionUsecase
}

func NewWalkOptionHandler(uc *usecase.WalkOptionUsecase) *WalkOptionHandler {
	return &WalkOptionHandler{uc: uc}
}

func (h *WalkOptionHandler) GetWalkOptions(c *echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.uc.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to get walk options",
		})
	}

	options := make([]openapi.WalkOptionResponse, 0, len(res.WalkOptions))
	for _, option := range res.WalkOptions {
		options = append(options, openapi.WalkOptionResponse{
			Id:       option.ID,
			Category: option.Category,
			Title:    option.Title,
		})
	}

	return c.JSON(http.StatusOK, openapi.WalkOptionListResponse{
		WalkOptions: options,
	})

}
