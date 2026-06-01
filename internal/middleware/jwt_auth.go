package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/auth"
	"github.com/labstack/echo/v5"
)

// func DummyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c *echo.Context) error {
// 		// 一旦userIDは固定。
// 		c.Set(
// 			"user_id",
// 			"019e1cd3-8194-7a36-816b-2f38206ca52c",
// 		)
// 		return next(c)
// 	}
// }

func JWTAuthMiddleware(jwtService *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				slog.WarnContext(c.Request().Context(), "missing authorization header")

				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "missing authorization header",
				})
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				slog.WarnContext(c.Request().Context(), "invalid authorization header")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid authorizaiton header",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				slog.ErrorContext(c.Request().Context(), "invalid token", "error", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid token",
				})
			}
			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}
}
