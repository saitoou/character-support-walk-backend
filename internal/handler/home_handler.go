package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chocoko/character-support-walk-backend/gen/openapi/v1"
	"github.com/chocoko/character-support-walk-backend/internal/usecase"
	"github.com/labstack/echo/v5"
)

type HomeHandler struct {
	uc             *usecase.HomeUsecase
	requestTimeout time.Duration
}

func NewHomeHandler(uc *usecase.HomeUsecase, requestTimeout time.Duration) *HomeHandler {
	return &HomeHandler{uc: uc, requestTimeout: requestTimeout}
}

func (h *HomeHandler) GetHome(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		h.requestTimeout,
	)
	defer cancel()

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "unauthorized",
		})
	}

	res, err := h.uc.GetHome(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusGatewayTimeout, openapi.ErrorResponse{
				Message: "request timeout",
			})
		}
		if errors.Is(err, usecase.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
				Message: "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to get home data",
		})
	}

	var todayWalk *openapi.HomeTodayWalk
	if res.TodayWalk != nil {
		todayWalk = &openapi.HomeTodayWalk{
			WalkId: res.TodayWalk.WalkID,
			Status: openapi.WalkStatus(res.TodayWalk.Status),
		}
	}
	return c.JSON(http.StatusOK, openapi.HomeResponse{
		Walker: openapi.HomeWalker{
			Nickname: res.Walker.Nickname,
		},
		Character: openapi.HomeCharacter{
			SupporterType: openapi.SupporterType(res.Character.SupporterType),
		},
		TodayWalk: todayWalk,
	})

}
