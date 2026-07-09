package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vhgomes/azgher/internal/api/dto"
	"github.com/vhgomes/azgher/internal/domain"
	"github.com/vhgomes/azgher/internal/repository"
	errPkg "github.com/vhgomes/azgher/pkg/errors"
	"github.com/vhgomes/azgher/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) error {
	logger.Info("verifying user email registration")

	_, err := s.repo.ByEmail(ctx, req.Email)
	if err == nil {
		logger.Info("email already registered", zap.String("email", req.Email))
		return errPkg.ErrEmailAlreadyRegistered
	}
	if !errors.Is(err, errPkg.ErrUserNotFound) {
		logger.Error("failed to verify email", err)
		return err
	}

	logger.Info("creating user", zap.String("name", req.Name))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", err)
		return err
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		logger.Error("failed to create user", err)
		return err
	}

	return nil
}

func (s *UserService) ById(ctx context.Context, id string) (*domain.User, error) {
	logger.Info("fetching user by id", zap.String("id", id))

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("failed to parse user id", err)
		return nil, err
	}

	user, err := s.repo.ByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			logger.Info("user not found", zap.String("id", id))
		} else {
			logger.Error("failed to fetch user by id", err)
		}
		return nil, err
	}

	logger.Info("user found", zap.String("user_id", user.ID.String()))
	return user, nil
}

func (s *UserService) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	logger.Info("fetching user by email")

	user, err := s.repo.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			logger.Info("user not found", zap.String("email", email))
		} else {
			logger.Error("failed to fetch user by email", err)
		}
		return nil, err
	}

	logger.Info("user found", zap.String("user_id", user.ID.String()))
	return user, nil
}

func (s *UserService) ByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	logger.Info("fetching user by google id", zap.String("google_id", googleID))

	user, err := s.repo.ByGoogleID(ctx, googleID)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			logger.Info("user not found", zap.String("google_id", googleID))
		} else {
			logger.Error("failed to fetch user by google id", err)
		}
		return nil, err
	}

	logger.Info("user found", zap.String("user_id", user.ID.String()))
	return user, nil
}

// TODO: precisa ser refatorado
func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	existingUser, err := s.repo.ByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			logger.Info("user not found", zap.String("id", user.ID.String()))
		} else {
			logger.Error("failed to verify user", err)
		}
		return err
	}

	if existingUser == nil {
		logger.Info("user not found", zap.String("id", user.ID.String()))
		return errPkg.ErrUserNotFound
	}

	// TODO: Temporario, ainda não sei se irei manter isso aqui
	updatedUser := &domain.User{
		ID:            existingUser.ID,
		Name:          existingUser.Name,
		Email:         existingUser.Email,
		PasswordHash:  existingUser.PasswordHash,
		GoogleID:      existingUser.GoogleID,
		AvatarURL:     existingUser.AvatarURL,
		EmailVerified: existingUser.EmailVerified,
		CreatedAt:     existingUser.CreatedAt,
		UpdatedAt:     existingUser.UpdatedAt,
		DeletedAt:     existingUser.DeletedAt,
	}

	// TODO: essas verificações são horripilantes, vai ser necessario modificar.
	if user.Name != "" {
		updatedUser.Name = user.Name
	}
	if user.Email != "" {
		updatedUser.Email = user.Email
	}
	if user.PasswordHash != "" {
		updatedUser.PasswordHash = user.PasswordHash
	}
	if user.GoogleID != "" {
		updatedUser.GoogleID = user.GoogleID
	}
	if user.AvatarURL != "" {
		updatedUser.AvatarURL = user.AvatarURL
	}

	_, err = s.repo.Update(ctx, updatedUser)
	if err != nil {
		logger.Error("failed to update user", err)
		return err
	}

	logger.Info("user updated", zap.String("id", updatedUser.ID.String()))
	return nil
}

func (s *UserService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	existingUser, err := s.repo.ByID(ctx, id)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			logger.Info("user not found", zap.String("id", id.String()))
		} else {
			logger.Error("failed to verify user", err)
		}
		return err
	}

	if existingUser == nil {
		logger.Info("user not found", zap.String("id", id.String()))
		return errPkg.ErrUserNotFound
	}

	err = s.repo.SoftDelete(ctx, id)
	if err != nil {
		logger.Error("failed to soft delete user", err)
		return err
	}

	logger.Info("user soft deleted", zap.String("id", id.String()))
	return nil
}
