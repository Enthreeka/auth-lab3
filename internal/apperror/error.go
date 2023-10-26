package apperror

import (
	"errors"
	"fmt"
)

var (
	ErrHashPasswordsNotEqual = NewAppError(errors.New("hashes_not_equal"), "Invalid password")
	ErrNoFoundRows           = NewAppError(errors.New("not_exist"), "Not found such data")
	ErrInvalidPassword       = NewAppError(errors.New("password_not_valid"), "Password must contain 8+symbols,at least one a-z and A-Z")
	ErrInvalidLogin          = NewAppError(errors.New("login_not_valid"), "Login must contain 6+symbol and 0-9A-Za-za-яА-Я")
)

type AppError struct {
	Err error  `json:"-"`
	Msg string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.Err)
}

func NewAppError(err error, msg string) *AppError {
	return &AppError{
		Err: err,
		Msg: msg,
	}
}
