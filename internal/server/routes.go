package server

import (
	"net/http"

	customMiddleware "template/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoutes() {
	e := s.Echo

	e.GET("/health", s.healthHandler)

	api := e.Group("/api/v1")

	// Auth Routes
	s.AuthHandler.RegisterRoutes(api)

	// Protected Routes
	protected := api.Group("")
	protected.Use(customMiddleware.Auth(s.Config.JWTSecret))
	s.UserHandler.RegisterRoutes(protected)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "up",
		"db":     s.DB.Health(),
		"redis":  s.Redis.Health(),
	})
}
