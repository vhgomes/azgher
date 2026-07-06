package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vhgomes/azgher/internal/domain"
	db "github.com/vhgomes/azgher/internal/postgres/db"
	pjerr "github.com/vhgomes/azgher/pkg/errors"
)

type UserRepo struct {
	queries *db.Queries
}

func NewUserRepo(queries *db.Queries) *UserRepo {
	return &UserRepo{queries: queries}
}

// ---------- Conversões ----------
func toDomain(u db.User) *domain.User {
	return &domain.User{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		PasswordHash:  stringOrEmpty(u.PasswordHash),
		GoogleID:      stringOrEmpty(u.GoogleID),
		AvatarURL:     stringOrEmpty(u.AvatarURL),
		Providers:     deriveProviders(u.PasswordHash, u.GoogleID),
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// deriveProviders calcula os providers com base nos campos preenchidos.
// Isso reflete exatamente o comentário no seu domain model:
// "Não persistido diretamente; calculado ou vindo de outra tabela/query"
func deriveProviders(passwordHash, googleID sql.NullString) []domain.AuthProvider {
	var providers []domain.AuthProvider
	if passwordHash.Valid && passwordHash.String != "" {
		providers = append(providers, domain.AuthProviderEmail)
	}
	if googleID.Valid && googleID.String != "" {
		providers = append(providers, domain.AuthProviderGoogle)
	}
	return providers
}

func stringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// ---------- Métodos ----------

func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: toNullString(user.PasswordHash),
		GoogleID:     toNullString(user.GoogleID),
		AvatarUrl:    toNullString(user.AvatarURL),
	})
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	domainUser := toDomain(row)
	return domainUser, nil
}

func (r *UserRepo) ByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", id, pjerr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("fetching user %s: %w", id, err)
	}
	return toDomain(row), nil
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s: %w", email, pjerr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("fetching user by email %s: %w", email, err)
	}
	return toDomain(row), nil
}

func (r *UserRepo) ByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	row, err := r.queries.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with google ID: %w", pjerr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("fetching user by google ID: %w", err)
	}
	return toDomain(row), nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		PasswordHash:  toNullString(user.PasswordHash),
		GoogleID:      toNullString(user.GoogleID),
		AvatarUrl:     toNullString(user.AvatarURL),
		EmailVerified: user.EmailVerified,
	})
	if err != nil {
		return nil, fmt.Errorf("updating user %s: %w", user.ID, err)
	}
	return toDomain(row), nil
}

func (r *UserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("soft deleting user %s: %w", id, err)
	}
	return nil
}
