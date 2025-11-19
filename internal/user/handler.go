package user

import (
	"errors"
	"net/http"

	"template/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(svc Service, jwtSecret string) *Handler {
	return &Handler{
		service:   svc,
		jwtSecret: jwtSecret,
	}
}

// RegisterRoutes now includes the /login route.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/register", h.handleRegister)
	g.POST("/login", h.handleLogin)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.GET("/me", h.handleMe)
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *Handler) handleRegister(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// NOTE: In a real application, you would add validation here using a library
	// like go-playground/validator to check the request struct tags.

	user, err := h.service.Register(c.Request().Context(), req.Username, req.Email, req.Password)
	if err != nil {
		// This is where you could check for specific domain errors
		// and return different status codes (e.g., 409 Conflict for duplicate email).
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not register user"})
	}

	return c.JSON(http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) handleLogin(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	user, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		// If the error is invalid credentials, return a 401 Unauthorized.
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		// For all other errors, return a 500.
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an unexpected error occurred"})
	}

	// Generate JWT access token
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
	}

	return c.JSON(http.StatusOK, loginResponse{
		AccessToken: accessToken,
	})
}

func (h *Handler) handleMe(c echo.Context) error {
	// The user's claims are automatically extracted by the middleware
	// and stored in the context.
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}

	// Now you have the user ID from the token, you can fetch the full user details.
	// We'll just return the claims for this example.
	// In a real app, you would call: user, err := h.service.GetUserByID(claims.UserID)
	return c.JSON(http.StatusOK, claims)
}
