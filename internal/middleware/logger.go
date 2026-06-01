package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c *echo.Context) bool {
			path := c.Request().URL.Path
			return path == "/health"
		},

		LogRequestID: true,
		LogMethod:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogStatus:    true,
		LogLatency:   true,

		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			userID, _ := c.Get("user_id").(string)

			attrs := []slog.Attr{
				slog.String("request_id", v.RequestID),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("user_id", userID),
			}

			if v.Error != nil {
				attrs = append(
					attrs,
					slog.String("error", v.Error.Error()),
				)
				slog.LogAttrs(
					c.Request().Context(),
					slog.LevelError,
					"request failed",
					attrs...,
				)
				return nil
			}
			slog.LogAttrs(
				c.Request().Context(),
				slog.LevelInfo,
				"request completed",
				attrs...,
			)
			return nil
		},
	})
}
