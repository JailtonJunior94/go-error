package http

import (
	"github.com/jailtonjunior94/go-error/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if domainErr, ok := err.(*domain.DomainError); ok {
			status := fiber.StatusInternalServerError
			switch domainErr.Code {
			case domain.ErrNotFound:
				status = fiber.StatusNotFound
			case domain.ErrInvalidInput:
				status = fiber.StatusBadRequest
			}

			// resposta JSON detalhada
			return c.Status(status).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    domainErr.Code,
					"message": domainErr.Message,
					"details": domainErr.Details,
				},
			})
		}

		if _, ok := err.(*domain.DomainError); !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    domain.ErrInternal,
					"message": "Erro inesperado",
				},
			})
		}

		return nil
	}
}
