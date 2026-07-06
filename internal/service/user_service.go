package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vhgomes/azgher/internal/domain"
	"github.com/vhgomes/azgher/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {

}

func (s *UserService) ById(ctx context.Context, id int) (*domain.User, error) {

}

func (s *UserService) ByEmail(ctx context.Context, email string) (*domain.User, error) {

}

func (s *UserService) ByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {

}

func (s *UserService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {

}

func (s *UserService) SoftDelete(ctx context.Context, id uuid.UUID) error {

}
