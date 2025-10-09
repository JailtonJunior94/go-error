package main

import (
	"context"
	"log"

	"github.com/jailtonjunior94/go-error/internal/application/usecase"
	infraHttp "github.com/jailtonjunior94/go-error/internal/infrastructure/http"
	"github.com/jailtonjunior94/go-error/pkg/o11y"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx := context.Background()

	metrics, shutdown, err := o11y.NewMetrics(ctx, "localhost:4317", "go-error", "1.0.0")
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown metrics: %v", err)
		}
	}()

	tracer, shutdown, err := o11y.NewTracer(ctx, "localhost:4317", "go-error", "1.0.0")
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	logger, shutdown, err := o11y.NewLogger(ctx, tracer, "localhost:4318", "go-error", "1.0.0")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown logger: %v", err)
		}
	}()

	telemetry, err := o11y.NewTelemetry(tracer, metrics, logger)
	if err != nil {
		log.Fatalf("failed to create telemetry: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: infraHttp.NewErrorHandler(),
	})

	userUseCase := usecase.NewUserUseCase(telemetry)
	userHandler := infraHttp.NewUserHandler(telemetry, userUseCase)

	app.Get("/users/:id", userHandler.GetUser)

	app.Listen(":8003")
}
