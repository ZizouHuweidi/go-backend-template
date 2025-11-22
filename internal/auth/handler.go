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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.RegisterRequest true "Register Request"
// @Success 201 {object} response.Response{data=jwt.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
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

// Login godoc
// @Summary Login user
// @Description Login with email and password to receive access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "Login Request"
// @Success 200 {object} response.Response{data=jwt.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Use a valid refresh token to get a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh Request"
// @Success 200 {object} response.Response{data=jwt.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/refresh [post]
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

// RecoverPassword godoc
// @Summary Request password recovery
// @Description Send a password recovery email to the user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RecoverPasswordRequest true "Recover Password Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/recover-password [post]
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

// ResetPassword godoc
// @Summary Reset password
// @Description Reset the user's password using a valid token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/reset-password [post]
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
