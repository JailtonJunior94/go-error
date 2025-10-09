package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jailtonjunior94/go-error/internal/domain"
	"github.com/jailtonjunior94/go-error/pkg/o11y"
)

type UserUseCase struct {
	telemetry o11y.Telemetry
}

func NewUserUseCase(telemetry o11y.Telemetry) *UserUseCase {
	return &UserUseCase{
		telemetry: telemetry,
	}
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (string, error) {
	ctx, span := uc.telemetry.Tracer().Start(ctx, "create_user_usecase.execute")
	defer span.End()

	start := time.Now()

	uc.telemetry.Metrics().AddCounter(ctx, "user_usecase_calls_total", 1, nil)
	uc.telemetry.Logger().Info(ctx, "GetUserByID called", o11y.Field{Key: "id", Value: id})

	if id == "" {
		span.SetAttributes(o11y.Attribute{Key: "error", Value: "ID do usuário não pode ser vazio"})
		return "", domain.NewDomainError(
			domain.ErrInvalidInput,
			"ID do usuário não pode ser vazio",
			map[string]any{"field": "id"},
			nil,
		)
	}

	if id == "500" {
		return "", fmt.Errorf("deu ruim")
	}

	if id != "123" {
		return "", domain.NewDomainError(
			domain.ErrNotFound,
			"Usuário não encontrado",
			map[string]any{"id": id},
			errors.New("registro ausente no banco"),
		)
	}

	uc.telemetry.Metrics().RecordHistogram(ctx, "user_fetch_duration_seconds", time.Since(start).Seconds(), nil)
	return "Usuário de exemplo", nil
}
