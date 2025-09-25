package domain

import "fmt"

type ErrorCode string

const (
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewDomainError(code ErrorCode, message string, details map[string]any, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}
