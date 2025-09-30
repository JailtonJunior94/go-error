package http

import (
	"github.com/jailtonjunior94/go-error/internal/application/usecase"
	"github.com/jailtonjunior94/go-error/pkg/o11y"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	telemetry o11y.Telemetry
	usecase   *usecase.UserUseCase
}

func NewUserHandler(
	telemetry o11y.Telemetry,
	usecase *usecase.UserUseCase,
) *UserHandler {
	return &UserHandler{
		telemetry: telemetry,
		usecase:   usecase,
	}
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
