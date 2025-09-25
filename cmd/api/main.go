package main

import (
	infraHttp "github.com/jailtonjunior94/go-error/internal/infrastructure/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jailtonjunior94/go-error/internal/application/usecase"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: infraHttp.NewErrorHandler(),
	})

	userUseCase := usecase.NewUserUseCase()
	userHandler := infraHttp.NewUserHandler(userUseCase)

	app.Get("/users/:id", userHandler.GetUser)

	app.Listen(":8003")

}
