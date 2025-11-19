package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("not found")
var ErrInvalidCredentials = errors.New("invalid email or password")

type Service interface {
	Register(ctx context.Context, username, email, password string) (User, error)
	Login(ctx context.Context, email, password string) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, username, email, password string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, username, email, string(hashedPassword))
	if err != nil {
		// TODO: Here you can check for specific DB errors, like a duplicate email,
		return User{}, err
	}

	return user, nil
}

// Login handles the user login business logic.
func (s *service) Login(ctx context.Context, email, password string) (User, error) {
	// Retrieve the user by email from the repository.
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// If the user is not found, return the invalid credentials error.
		// We don't want to tell the attacker whether the email exists.
		return User{}, ErrInvalidCredentials
	}

	// Compare the provided password with the stored hash.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// If the password does not match, return the same invalid credentials error.
		return User{}, ErrInvalidCredentials
	}

	// Login successful.
	return user, nil
}
