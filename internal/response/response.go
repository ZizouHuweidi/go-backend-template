package response

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Error struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func JSON(c echo.Context, status int, data interface{}, meta interface{}) error {
	return c.JSON(status, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorJSON(c echo.Context, status int, code, message string, details interface{}) error {
	return c.JSON(status, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
