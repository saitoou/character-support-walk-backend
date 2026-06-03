package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/chocoko/character-support-walk-backend/gen/openapi/v1"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/auth"
	"github.com/chocoko/character-support-walk-backend/internal/usecase"
	"github.com/labstack/echo/v5"
)

type AuthUsecase interface {
	LoginWithGoogle(ctx context.Context, idToken string) (usecase.UserTokens, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (usecase.RefreshOutput, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthHandler struct {
	uc         AuthUsecase
	jwtService *auth.JWTService
}

func NewAuthHandler(uc AuthUsecase, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		uc:         uc,
		jwtService: jwtService,
	}
}

// Test Login func
func (h *AuthHandler) DevLogin(c *echo.Context) error {
	token, err := h.jwtService.GenerateAccessToken(
		"019e1cd3-8194-7a36-816b-2f38206ca52c",
	)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, openapi.DevLoginResponse{
		AccessToken: token,
	})
}

func (h *AuthHandler) GoogleLogin(c *echo.Context) error {
	ctx := c.Request().Context()

	var req openapi.LoginWithGoogleJSONRequestBody

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "invalid request",
		})
	}

	res, err := h.uc.LoginWithGoogle(ctx, req.IdToken)
	if err != nil {
		slog.Error("failed to get user info", "error", err)
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "failed to Login with Google",
		})
	}

	return c.JSON(http.StatusOK, openapi.AuthTokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})

}

func (h *AuthHandler) Refresh(c *echo.Context) error {
	ctx := c.Request().Context()

	var req openapi.RefreshAccessTokenJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "invalid request",
		})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "refresh_token is required",
		})
	}

	res, err := h.uc.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "failed to refresh access token",
		})
	}

	return c.JSON(http.StatusOK, openapi.RefreshTokenResponse{
		AccessToken: res.AccessToken,
	})
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	ctx := c.Request().Context()
	var req openapi.LogoutRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request"})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "refresh_token is required",
		})
	}

	if err := h.uc.Logout(ctx, req.RefreshToken); err != nil {
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to logout"})
	}

	return c.NoContent(http.StatusNoContent)

}
