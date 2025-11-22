package user

import (
	"context"
	"errors"
	"fmt"
	"template/internal/email"
	"template/internal/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*jwt.TokenPair, error)
	Login(ctx context.Context, req *LoginRequest) (*jwt.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type service struct {
	repo         Repository
	jwtSecret    string
	emailSender  *email.Sender
	frontendHost string
}

func NewService(repo Repository, jwtSecret string, emailSender *email.Sender, frontendHost string) Service {
	return &service{
		repo:         repo,
		jwtSecret:    jwtSecret,
		emailSender:  emailSender,
		frontendHost: frontendHost,
	}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*jwt.TokenPair, error) {
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*jwt.TokenPair, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *service) RefreshToken(ctx context.Context, token string) (*jwt.TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, ErrInvalidToken
	}

	// Reuse Detection
	if rt.Revoked {
		// Token reused! Revoke all tokens for this user (Family Tracking)
		_ = s.repo.RevokeAllUserTokens(ctx, rt.UserID)
		return nil, ErrInvalidToken
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Revoke the used refresh token (Rotation)
	err = s.repo.RevokeRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, rt.UserID)
}

func (s *service) generateTokens(ctx context.Context, userID string) (*jwt.TokenPair, error) {
	tokens, err := jwt.GenerateTokens(userID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken := &RefreshToken{
		UserID:    userID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	err = s.repo.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	// Generate a short-lived token
	token, err := jwt.GenerateResetToken(user.ID, s.jwtSecret)
	if err != nil {
		return err
	}

	// Send email
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendHost, token)
	body := fmt.Sprintf("Click here to reset your password: <a href=\"%s\">Reset Password</a>", resetLink)

	return s.emailSender.Send(email, "Password Recovery", body)
}

func (s *service) ResetPassword(ctx context.Context, tokenString, newPassword string) error {
	claims, err := jwt.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return ErrInvalidToken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, claims.UserID, string(hashedPassword))
}
