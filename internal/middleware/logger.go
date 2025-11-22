package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func SlogLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			latency := time.Since(start)
			status := c.Response().Status
			method := c.Request().Method
			path := c.Request().URL.Path
			ip := c.RealIP()

			msg := "incoming request"
			attrs := []slog.Attr{
				slog.Int("status", status),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", ip),
				slog.Duration("latency", latency),
			}

			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
				logger.LogAttrs(c.Request().Context(), slog.LevelError, msg, attrs...)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, msg, attrs...)
			}

			return err
		}
	}
}
