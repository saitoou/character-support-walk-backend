package middleware

import (
	"log/slog"

	"github.com/chocoko/character-support-walk-backend/internal/apperr"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c *echo.Context) bool {
			path := c.Request().URL.Path
			return path == "/health" || path == "/mertrics"
		},

		LogRequestID: true,
		LogMethod:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogStatus:    true,
		LogLatency:   true,

		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			userID, _ := c.Get("user_id").(string)

			span := trace.SpanFromContext(c.Request().Context())
			spanCtx := span.SpanContext()

			traceID := ""
			spanID := ""

			if spanCtx.IsValid() {
				traceID = spanCtx.TraceID().String()
				spanID = spanCtx.SpanID().String()
			}

			attrs := []slog.Attr{
				slog.String("request_id", v.RequestID),
				slog.String("trace_id", traceID),
				slog.String("span_id", spanID),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("user_id", userID),
			}

			if v.Error != nil {
				if appErr, ok := apperr.AsAppError(v.Error); ok {
					attrs = append(
						attrs,
						slog.String("error_message", appErr.Message),
						slog.String("error", v.Error.Error()),
						slog.String("error_file", appErr.File),
						slog.String("error_func", appErr.Func),
						slog.Int("error_line", appErr.Line),
					)
				} else {
					attrs = append(attrs,
						slog.String("error", v.Error.Error()))
				}
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
