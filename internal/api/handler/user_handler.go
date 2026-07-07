package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vhgomes/azgher/internal/api/dto"
	"github.com/vhgomes/azgher/internal/service"
	"github.com/vhgomes/azgher/pkg/validator"
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

	err := h.UserService.Create(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error creating user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
}
