package usecase

import (
	"errors"
	"fmt"

	"github.com/jailtonjunior94/go-error/internal/domain"
)

type UserUseCase struct{}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{}
}

func (uc *UserUseCase) GetUserByID(id string) (string, error) {
	if id == "" {
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

	return "Usuário de exemplo", nil
}
