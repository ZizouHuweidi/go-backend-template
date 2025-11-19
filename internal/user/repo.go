package user

import (
	"context"

	"template/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type postgresRepository struct {
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{
		queries: db.New(pool),
	}
}

func (r *postgresRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (User, error) {
	createdUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return User{}, err
	}

	return User{
		ID:           createdUser.ID,
		Username:     createdUser.Username,
		Email:        createdUser.Email,
		PasswordHash: createdUser.PasswordHash,
		CreatedAt:    createdUser.CreatedAt.Time,
		UpdatedAt:    createdUser.UpdatedAt.Time,
	}, nil
}

func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}, nil
}
