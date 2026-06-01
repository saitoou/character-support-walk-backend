package handler

import (
	"errors"
	"net/http"

	"github.com/chocoko/character-support-walk-backend/gen/openapi/v1"
	"github.com/chocoko/character-support-walk-backend/internal/usecase"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	uc *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) GetUser(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	res, err := h.uc.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "failed to get user info",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to get walker data",
		})
	}

	return c.JSON(http.StatusOK, openapi.WalkerResponse{
		Id:            res.UserID,
		Nickname:      res.Nickname,
		SupporterType: openapi.SupporterType(res.SupporterType),
		CreatedAt:     res.CreatedAt,
		UpdatedAt:     res.UpdatedAt,
	})
}

func (h *UserHandler) UpdateWalker(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	var req openapi.UpdateWalkerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "invalid request",
		})
	}
	res, err := h.uc.UpdateByID(ctx, userID, req.Nickname, string(req.SupporterType))
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to update user profile",
		})
	}

	return c.JSON(http.StatusOK, openapi.WalkerResponse{
		Id:            res.UserID,
		Nickname:      res.Nickname,
		SupporterType: openapi.SupporterType(res.SupporterType),
		CreatedAt:     res.CreatedAt,
		UpdatedAt:     res.UpdatedAt,
	})
}

func (h *UserHandler) DeactivateWalker(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}
	err := h.uc.Deactivate(ctx, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to delete user profile",
		})
	}
	return c.NoContent(http.StatusNoContent)
}
