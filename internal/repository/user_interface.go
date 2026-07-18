package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vhgomes/azgher/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	ById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
	ByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit int, offset int) ([]*domain.User, error)
}
