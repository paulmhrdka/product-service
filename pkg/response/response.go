package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status     int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       interface{}    `json:"data,omitempty"`
	Error      interface{}    `json:"error,omitempty"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PageSize    int   `json:"page_size"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// Success response without pagination
func Success(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Success response with pagination
func SuccessWithPagination(c echo.Context, data interface{}, pagination PaginationMeta, message string) error {
	return c.JSON(http.StatusOK, Response{
		Status:     http.StatusOK,
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// Bad Request error
func BadRequest(c echo.Context, err interface{}, message string) error {
	return c.JSON(http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Not Found error
func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, Response{
		Status:  http.StatusNotFound,
		Success: false,
		Message: message,
	})
}

// Internal Server error
func InternalServerError(c echo.Context, err interface{}) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Status:  http.StatusInternalServerError,
		Success: false,
		Message: "Internal server error",
		Error:   err,
	})
}

// Validation error
func ValidationError(c echo.Context, errors []ErrorDetail) error {
	return c.JSON(http.StatusUnprocessableEntity, Response{
		Status:  http.StatusUnprocessableEntity,
		Success: false,
		Message: "Validation failed",
		Error:   errors,
	})
}
