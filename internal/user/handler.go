package user

import (
	"net/http"

	"template/internal/json"
	"template/internal/jwt"
	"template/internal/response"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/users/me", h.GetMe)
}

func (h *Handler) GetMe(c echo.Context) error {
	claims, ok := c.Get("user").(*jwt.Claims)
	if !ok {
		return json.Unauthorized(c, "Invalid token")
	}

	user, err := h.repo.GetByID(c.Request().Context(), claims.UserID)
	if err != nil {
		return json.InternalServerError(c, err)
	}
	if user == nil {
		return json.NotFound(c, "User not found")
	}

	return response.JSON(c, http.StatusOK, user, nil)
}
