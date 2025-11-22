package auth

import (
	"net/http"

	"template/internal/json"
	"template/internal/response"
	"template/internal/user"
	"template/internal/validator"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	userService user.Service
	validator   *validator.Validator
}

func NewHandler(userService user.Service, validator *validator.Validator) *Handler {
	return &Handler{
		userService: userService,
		validator:   validator,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/auth/register", h.Register)
	g.POST("/auth/login", h.Login)
	g.POST("/auth/refresh", h.RefreshToken)
	g.POST("/auth/recover-password", h.RecoverPassword)
	g.POST("/auth/reset-password", h.ResetPassword)
}

func (h *Handler) Register(c echo.Context) error {
	var req user.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return json.BadRequest(c, err)
	}

	if err := h.validator.Validate(req); err != nil {
		return json.BadRequest(c, err)
	}

	tokens, err := h.userService.Register(c.Request().Context(), &req)
	if err != nil {
		if err == user.ErrUserAlreadyExists {
			return response.ErrorJSON(c, http.StatusConflict, "USER_ALREADY_EXISTS", "User with this email already exists", nil)
		}
		return json.InternalServerError(c, err)
	}

	return response.JSON(c, http.StatusCreated, tokens, nil)
}

func (h *Handler) Login(c echo.Context) error {
	var req user.LoginRequest
	if err := c.Bind(&req); err != nil {
		return json.BadRequest(c, err)
	}

	if err := h.validator.Validate(req); err != nil {
		return json.BadRequest(c, err)
	}

	tokens, err := h.userService.Login(c.Request().Context(), &req)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			return json.Unauthorized(c, "Invalid credentials")
		}
		return json.InternalServerError(c, err)
	}

	return response.JSON(c, http.StatusOK, tokens, nil)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return json.BadRequest(c, err)
	}

	if err := h.validator.Validate(req); err != nil {
		return json.BadRequest(c, err)
	}

	tokens, err := h.userService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if err == user.ErrInvalidToken {
			return json.Unauthorized(c, "Invalid or expired refresh token")
		}
		return json.InternalServerError(c, err)
	}

	return response.JSON(c, http.StatusOK, tokens, nil)
}

type RecoverPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *Handler) RecoverPassword(c echo.Context) error {
	var req RecoverPasswordRequest
	if err := c.Bind(&req); err != nil {
		return json.BadRequest(c, err)
	}

	if err := h.validator.Validate(req); err != nil {
		return json.BadRequest(c, err)
	}

	err := h.userService.ForgotPassword(c.Request().Context(), req.Email)
	if err != nil {
		return json.InternalServerError(c, err)
	}

	return response.JSON(c, http.StatusOK, map[string]string{"message": "If the email exists, a recovery link has been sent."}, nil)
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return json.BadRequest(c, err)
	}

	if err := h.validator.Validate(req); err != nil {
		return json.BadRequest(c, err)
	}

	err := h.userService.ResetPassword(c.Request().Context(), req.Token, req.NewPassword)
	if err != nil {
		if err == user.ErrInvalidToken {
			return json.Unauthorized(c, "Invalid or expired token")
		}
		return json.InternalServerError(c, err)
	}

	return response.JSON(c, http.StatusOK, map[string]string{"message": "Password updated successfully"}, nil)
}
