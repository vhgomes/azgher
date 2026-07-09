package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vhgomes/azgher/internal/domain"
	db "github.com/vhgomes/azgher/internal/postgres/db"
	pgutil "github.com/vhgomes/azgher/internal/postgres/pgutil"
	errPkg "github.com/vhgomes/azgher/pkg/errors"
	"github.com/vhgomes/azgher/pkg/logger"
	"go.uber.org/zap"
)

type UserRepo struct {
	queries *db.Queries
}

func NewUserRepo(queries *db.Queries) *UserRepo {
	return &UserRepo{queries: queries}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: pgutil.NilIfEmpty(user.PasswordHash),
		GoogleID:     pgutil.NilIfEmpty(user.GoogleID),
		AvatarUrl:    pgutil.NilIfEmpty(user.AvatarURL),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.Info("email already registered", zap.String("email", user.Email))
			return errPkg.ErrEmailAlreadyRegistered
		}
		logger.Error("failed to create user", err, zap.String("email", user.Email))
		return err
	}

	return nil
}

func (r *UserRepo) ByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("user not found", zap.String("user_id", id.String()))
			return nil, fmt.Errorf("user %s: %w", id, errPkg.ErrUserNotFound)
		}
		logger.Error("failed to fetch user by id", err, zap.String("user_id", id.String()))
		return nil, fmt.Errorf("fetching user %s: %w", id, err)
	}
	return domain.NewUser(
		row.ID, row.Name, row.Email,
		row.PasswordHash, row.GoogleID, row.AvatarUrl,
		row.EmailVerified,
		row.CreatedAt, row.UpdatedAt,
	), nil
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("user not found by email")
			return nil, fmt.Errorf("user with email: %w", errPkg.ErrUserNotFound)
		}
		logger.Error("failed to fetch user by email", err, zap.String("email", email))
		return nil, fmt.Errorf("fetching user by email: %w", err)
	}
	return domain.NewUser(
		row.ID, row.Name, row.Email,
		row.PasswordHash, row.GoogleID, row.AvatarUrl,
		row.EmailVerified,
		row.CreatedAt, row.UpdatedAt,
	), nil
}

func (r *UserRepo) ByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	row, err := r.queries.GetUserByGoogleID(ctx, &googleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("user not found by google id")
			return nil, fmt.Errorf("user with google ID: %w", errPkg.ErrUserNotFound)
		}
		logger.Error("failed to fetch user by google id", err)
		return nil, fmt.Errorf("fetching user by google ID: %w", err)
	}
	return domain.NewUser(
		row.ID, row.Name, row.Email,
		row.PasswordHash, row.GoogleID, row.AvatarUrl,
		row.EmailVerified,
		row.CreatedAt, row.UpdatedAt,
	), nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		PasswordHash:  pgutil.NilIfEmpty(user.PasswordHash),
		GoogleID:      pgutil.NilIfEmpty(user.GoogleID),
		AvatarUrl:     pgutil.NilIfEmpty(user.AvatarURL),
		EmailVerified: user.EmailVerified,
	})
	if err != nil {
		logger.Error("failed to update user", err, zap.String("user_id", user.ID.String()))
		return nil, fmt.Errorf("updating user %s: %w", user.ID, err)
	}

	logger.Info("user updated", zap.String("user_id", row.ID.String()))

	return domain.NewUser(
		row.ID, row.Name, row.Email,
		row.PasswordHash, row.GoogleID, row.AvatarUrl,
		row.EmailVerified,
		row.CreatedAt, row.UpdatedAt,
	), nil
}

func (r *UserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.SoftDeleteUser(ctx, id)
	if err != nil {
		logger.Error("failed to soft delete user", err, zap.String("user_id", id.String()))
		return fmt.Errorf("soft deleting user %s: %w", id, err)
	}

	logger.Info("user soft deleted", zap.String("user_id", id.String()))
	return nil
}

func (r *UserRepo) List(ctx context.Context, limit, offset int32) ([]domain.User, error) {
	rows, err := r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to list users", err, zap.Int32("limit", limit), zap.Int32("offset", offset))
		return nil, fmt.Errorf("listing users: %w", err)
	}

	users := make([]domain.User, len(rows))
	for i, row := range rows {
		u := domain.NewUser(
			row.ID, row.Name, row.Email,
			row.PasswordHash, row.GoogleID, row.AvatarUrl,
			row.EmailVerified,
			row.CreatedAt, row.UpdatedAt,
		)
		users[i] = *u
	}

	logger.Info("users listed", zap.Int("count", len(users)), zap.Int32("limit", limit), zap.Int32("offset", offset))
	return users, nil
}
