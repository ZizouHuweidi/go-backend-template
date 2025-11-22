package user

import (
	"time"
)

type User struct {
	ID           string     `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	Username     string     `db:"username" json:"username"`
	PasswordHash string     `db:"password_hash" json:"-"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	LastLogin    *time.Time `db:"last_login" json:"last_login,omitempty"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	Revoked   bool      `db:"revoked"`
}
