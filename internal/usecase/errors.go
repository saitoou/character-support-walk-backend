package usecase

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternalServer = errors.New("internal server error")
)
