package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}

	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// error implements the fmt.Error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, message string) Error {
	return Error{
		Code: code,
		Err:  message,
	}
}

func ErrorInvalidId() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id",
	}
}

func ErrorUnathorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unathorized",
	}
}

func ErrorBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}
}

func ErrorResourceNotFound() Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  "resource not found",
	}
}
