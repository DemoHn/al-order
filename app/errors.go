package app

import (
	"fmt"

	echo "github.com/labstack/echo"
)

// Error - customized AppError for displaying error response
type Error struct {
	Message    string
	StatusCode int
	Code       int
	Details    string
}

// ErrIDNotFound - if input ID is empty
func ErrIDNotFound(details string) *Error {
	return &Error{
		Message:    "orderID Not Found",
		StatusCode: 400,
		Code:       10000,
		Details:    details,
	}
}

// ErrDBFatal - db operation error
func ErrDBFatal(err error) *Error {
	return &Error{
		Message:    "fatal error: db operation error",
		StatusCode: 500,
		Code:       10001,
		Details:    err.Error(),
	}
}

// ErrRedisFatal - redis operation error
func ErrRedisFatal(err error) *Error {
	return &Error{
		Message:    "fatal error: redis operation error",
		StatusCode: 500,
		Code:       10002,
		Details:    err.Error(),
	}
}

// ErrGoogleMapService - googleMap service error
func ErrGoogleMapService(err error) *Error {
	return &Error{
		Message:    "fetch google map direction data failed",
		StatusCode: 500,
		Code:       10003,
		Details:    err.Error(),
	}
}

// ErrNoRoute - if there's no available route from origin -> dest
func ErrNoRoute() *Error {
	return &Error{
		Message:    "there's no available route from origin -> destination",
		StatusCode: 400,
		Code:       10004,
		Details:    "",
	}
}

// ErrOrderHasTaken -
func ErrOrderHasTaken() *Error {
	return &Error{
		Message:    "order has been taken by somebody else",
		StatusCode: 400,
		Code:       10005,
		Details:    "",
	}
}

// ErrUnknownFatal - unknown fatal error
func ErrUnknownFatal(err error) *Error {
	return &Error{
		Message:    "unknown fatal error",
		StatusCode: 500,
		Code:       11000,
		Details:    err.Error(),
	}
}

// ErrGeneralHTTP - from echo's internal server
func ErrGeneralHTTP(err *echo.HTTPError) *Error {
	return &Error{
		Message:    fmt.Sprintf("general HTTP error: %s", err.Error()),
		StatusCode: err.Code,
		Code:       12000,
		Details:    fmt.Sprintf("%v", err.Message),
	}
}

// implement type `error` interface:
func (e *Error) Error() string {
	return e.Message
}
