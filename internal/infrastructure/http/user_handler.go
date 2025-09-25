package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jailtonjunior94/go-error/internal/application/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUseCase
}

func NewUserHandler(u *usecase.UserUseCase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		// apenas retorna o erro, middleware tratará
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": user,
	})
}
