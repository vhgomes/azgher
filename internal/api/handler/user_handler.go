package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/vhgomes/azgher/internal/api/dto"
	"github.com/vhgomes/azgher/internal/service"
	errPkg "github.com/vhgomes/azgher/pkg/errors"
	"github.com/vhgomes/azgher/pkg/logger"
	"github.com/vhgomes/azgher/pkg/validator"
	"go.uber.org/zap"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	if validationErrs := validator.Validate(&req); validationErrs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": validationErrs,
		})
	}

	if err := h.UserService.Create(c.Context(), req); err != nil {
		if errors.Is(err, errPkg.ErrEmailAlreadyRegistered) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "email already registered",
			})
		}
		logger.Error("failed to create user", err, zap.String("email", req.Email))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error creating user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
}

func (h *UserHandler) GetById(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	idUser, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid format",
		})
	}

	user, err := h.UserService.ById(c.Context(), idUser)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		logger.Error("failed to fetch user by id", err, zap.String("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error fetching user by id",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewUserResponse(user))
}

func (h *UserHandler) GetByEmail(c fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}

	user, err := h.UserService.ByEmail(c.Context(), email)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		logger.Error("failed to fetch user by email", err, zap.String("email", email))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error fetching user by email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewUserResponse(user))
}

func (h *UserHandler) ByGoogleId(c fiber.Ctx) error {
	googleId := c.Params("google_id")
	if googleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "google_id is required",
		})
	}

	user, err := h.UserService.ByGoogleID(c.Context(), googleId)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		logger.Error("failed to fetch user by google_id", err, zap.String("google_id", googleId))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewUserResponse(user))
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	if validationErrs := validator.Validate(&req); validationErrs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": validationErrs,
		})
	}

	err := h.UserService.Update(c.Context(), req)
	if errors.Is(err, errPkg.ErrUserNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	if errors.Is(err, errPkg.ErrEmailAlreadyRegistered) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.Status(fiber.StatusOK).JSON("user updated successfully")
}
