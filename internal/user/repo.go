package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
}

type repository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query, args, err := r.sb.Insert("users").
		Columns("email", "username", "password_hash").
		Values(user.Email, user.Username, user.PasswordHash).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	return r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query, args, err := r.sb.Select("*").From("users").Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	query, args, err := r.sb.Select("*").From("users").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	query, args, err := r.sb.Insert("refresh_tokens").
		Columns("user_id", "token", "expires_at").
		Values(token.UserID, token.Token, token.ExpiresAt).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	return r.db.QueryRowContext(ctx, query, args...).Scan(&token.ID, &token.CreatedAt)
}

func (r *repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	query, args, err := r.sb.Select("*").From("refresh_tokens").Where(squirrel.Eq{"token": token}).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &rt, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, token string) error {
	query, args, err := r.sb.Update("refresh_tokens").
		Set("revoked", true).
		Where(squirrel.Eq{"token": token}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	query, args, err := r.sb.Update("refresh_tokens").
		Set("revoked", true).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *repository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query, args, err := r.sb.Update("users").
		Set("password_hash", passwordHash).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
