package middleware

import (
	"strings"

	"template/internal/json"
	"template/internal/jwt"

	"github.com/labstack/echo/v4"
)

func Auth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return json.Unauthorized(c, "Missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return json.Unauthorized(c, "Invalid authorization header format")
			}

			tokenString := parts[1]
			claims, err := jwt.ValidateToken(tokenString, secret)
			if err != nil {
				return json.Unauthorized(c, "Invalid or expired token")
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
