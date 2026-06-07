package apperr

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

type AppError struct {
	Message string
	Err     error
	File    string
	Func    string
	Line    int
}

func Wrap(message string, err error) *AppError {
	if err == nil {
		return nil
	}

	pc, file, line, ok := runtime.Caller(1)

	funcName := ""
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
		}
	}

	return &AppError{
		Message: message,
		Err:     err,
		File:    filepath.Base(file),
		Func:    funcName,
		Line:    line,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
