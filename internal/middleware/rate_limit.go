package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(1000),
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "too many request",
			})
		},
	})
}
