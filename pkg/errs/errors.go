package errs

import (
	"errors"
	"net/http"
)

type Error struct {
	code    int
	message string
	err     error
}

func (e *Error) Error() string {
	return e.message
}

func NewWithMessage(code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
		err:     errors.New(message),
	}
}

func GetHTTPCode(err error) int {
	var errWrapped *Error
	ok := errors.As(err, &errWrapped)
	if ok {
		return errWrapped.HTTPCode()
	}
	return http.StatusInternalServerError
}

func (e *Error) HTTPCode() int {
	return e.code
}

// New error object
func Wrap(code int, err error) *Error {
	return &Error{
		code:    code,
		message: err.Error(),
		err:     err,
	}
}
