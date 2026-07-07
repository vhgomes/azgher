package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vhgomes/azgher/internal/api/dto"
	"github.com/vhgomes/azgher/internal/domain"
	"github.com/vhgomes/azgher/internal/repository"
	"github.com/vhgomes/azgher/pkg/logger"
	"go.uber.org/zap"

	"github.com/vhgomes/azgher/pkg/errors"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, dto dto.CreateUserRequest) error {
	logger.Info("verifing if %s is already registered", zap.String("email", dto.Email))
	userExists, err := s.repo.ByEmail(ctx, dto.Email)

	if err != nil {
		logger.Error("failed to verify email", err)
		return err
	}

	if userExists != nil {
		logger.Info("email already registered")
		return errors.ErrEmailAlreadyRegistered
	}

	user := dto.ToDomain()

	logger.Info("creating user", zap.String("name", user.Name), zap.String("email", user.Email))
	err = s.repo.Create(ctx, user)

	if err != nil {
		logger.Error("failed to create user", err)
		return err
	}

	return nil
}

// TODO: Retornar DTO
func (s *UserService) ById(ctx context.Context, id int) error {

}

// TODO: Retornar DTO
func (s *UserService) ByEmail(ctx context.Context, email string) error {

}

// TODO: Retornar DTO
func (s *UserService) ByGoogleID(ctx context.Context, googleID string) error {

}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {

}

func (s *UserService) SoftDelete(ctx context.Context, id uuid.UUID) error {

}
