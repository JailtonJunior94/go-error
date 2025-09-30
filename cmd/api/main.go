package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/jailtonjunior94/go-error/internal/application/usecase"
	infraHttp "github.com/jailtonjunior94/go-error/internal/infrastructure/http"
	"github.com/jailtonjunior94/go-error/pkg/o11y"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx := context.Background()

	metrics, shutdown, err := o11y.NewMetrics(ctx, "localhost:4317", "go-error", nil)
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown metrics: %v", err)
		}
	}()

	tracer, shutdown, err := o11y.NewTracer(ctx, "localhost:4317", "go-error", nil)
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	logger := o11y.NewLogger(slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{})))

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
