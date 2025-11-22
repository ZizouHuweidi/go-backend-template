package middleware

import (
	"fmt"
	"net/http"
	"time"

	"template/internal/redis"
	"template/internal/response"

	"github.com/labstack/echo/v4"
)

func RateLimit(redisClient *redis.Client, limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := fmt.Sprintf("rate_limit:%s", ip)

			count, err := redisClient.Client.Incr(c.Request().Context(), key).Result()
			if err != nil {
				return next(c) // Fail open if Redis is down
			}

			if count == 1 {
				redisClient.Client.Expire(c.Request().Context(), key, window)
			}

			if count > int64(limit) {
				return response.ErrorJSON(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests", nil)
			}

			return next(c)
		}
	}
}
