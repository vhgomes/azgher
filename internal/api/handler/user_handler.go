package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vhgomes/azgher/internal/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) Create(fiberCtx *fiber.Ctx) {

}
