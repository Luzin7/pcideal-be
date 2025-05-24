package errors

import (
	"fmt"
)

type ErrService struct {
	StatusCode int
	Message    string
}

func (e *ErrService) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func New(msg string, statusCode int) *ErrService {
	return &ErrService{
		Message:    msg,
		StatusCode: statusCode,
	}
}

func ErrNotFound(ctx string) *ErrService {
	return New(
		fmt.Sprintf("%s not found", ctx),
		404,
	)
}

func ErrAlreadyExists(ctx string) *ErrService {
	return New(
		fmt.Sprintf("%s already exists", ctx),
		409,
	)
}

func ErrInternalServerError() *ErrService {
	return New(
		"Internal server error",
		500,
	)
}
