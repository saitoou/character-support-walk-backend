package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/chocoko/character-support-walk-backend/gen/openapi/v1"
	"github.com/chocoko/character-support-walk-backend/internal/usecase"
	"github.com/labstack/echo/v5"
)

type WalkUsecase interface {
	Start(ctx context.Context, walkOptionId int, userID string) (usecase.StartWalkOutput, error)
	GetWalks(ctx context.Context, userID string) (usecase.WalkListOutput, error)
	GetWalking(ctx context.Context, userID string, walkID string) (usecase.WalkOutput, error)
	UpdateComplete(ctx context.Context, userID, walkID string) error
	UpdateCancel(ctx context.Context, userID, walkID string) error
}

type WalkHandler struct {
	uc *usecase.WalkUsecase
}

func NewWalkHandler(uc *usecase.WalkUsecase) *WalkHandler {
	return &WalkHandler{uc: uc}
}

func (h *WalkHandler) StartWalk(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	var req openapi.StartWalkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "invalid request",
		})
	}
	res, err := h.uc.Start(ctx, req.WalkOptionId, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to start walk",
		})
	}

	return c.JSON(http.StatusOK, openapi.StartWalkResponse{
		WalkId: res.WalkID,
		Status: openapi.WalkStatus(res.Status),
	})

}

func (h *WalkHandler) GetWalks(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	walks, err := h.uc.GetWalks(ctx, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "failed to get walks",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to get walks data",
		})
	}

	res := make([]openapi.WalkRecordResponse, 0, len(walks.Walks))
	for _, walk := range walks.Walks {
		res = append(res, openapi.WalkRecordResponse{
			Id:           walk.ID,
			WalkOptionId: walk.WalkOptionID,
			Status:       openapi.WalkStatus(walk.Status),
			StartedAt:    walk.StartedAt,
			FinishedAt:   walk.FinishedAt,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *WalkHandler) GetWalk(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}
	walkID := c.Param("walkID")

	walk, err := h.uc.GetWalking(ctx, userID, walkID)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to get walks",
		})
	}

	return c.JSON(http.StatusOK, openapi.WalkRecordResponse{
		Id:           walk.ID,
		WalkOptionId: walk.WalkOptionID,
		Status:       openapi.WalkStatus(walk.Status),
		StartedAt:    walk.StartedAt,
		FinishedAt:   walk.FinishedAt,
	})
}

func (h *WalkHandler) CompleteWalk(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	walkID := c.Param("walkID")

	err := h.uc.UpdateComplete(ctx, userID, walkID)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to update walks",
		})
	}

	return c.JSON(http.StatusOK, openapi.MessageResponse{
		Message: "walk completed",
	})
}

func (h *WalkHandler) CancelWalk(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	walkID := c.Param("walkID")

	err := h.uc.UpdateCancel(ctx, userID, walkID)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to cancel walks",
		})
	}

	return c.JSON(http.StatusOK, openapi.MessageResponse{
		Message: "walk canceled",
	})
}
