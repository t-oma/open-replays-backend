// Package response provides HTTP response helpers.
package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"open-replays/internal/api/httperr"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    httperr.ErrorCode `json:"code"`
	Message string            `json:"message"`
	Details any               `json:"details,omitempty"`
}

// SuccessResponse represents a successful response.
type SuccessResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
}

type Pagination struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int   `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

func NewPagination(page, pageSize int, totalItems int64) Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return Pagination{
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}
}

type PaginatedData[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// OK sends a successful response with data.
func OK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, SuccessResponse[T]{Data: data})
}

// OKWithMessage sends a successful response with data and message.
func OKWithMessage[T any](c *gin.Context, data T, message string) {
	c.JSON(http.StatusOK, SuccessResponse[T]{Data: data, Message: message})
}

// Created sends a 201 Created response.
func Created[T any](c *gin.Context, data T) {
	c.JSON(http.StatusCreated, SuccessResponse[T]{Data: data})
}

// NoContent sends a 204 No Content response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response.
func Error(c *gin.Context, err *httperr.AppError) {
	c.JSON(err.HTTPStatus, ErrorResponse{
		Code:    err.Code,
		Message: err.Message,
		Details: err.Details,
	})
}

// ErrorFromError maps any error to HTTP response.
func ErrorFromError(c *gin.Context, err error) {
	appErr := httperr.MapError(err)
	Error(c, appErr)
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(c *gin.Context, message string, details ...any) {
	err := httperr.New(httperr.ErrCodeInvalidRequest, message, http.StatusBadRequest)
	if len(details) > 0 {
		err = err.WithDetails(details[0])
	}
	Error(c, err)
}

// ValidationFailed sends a 400 Bad Request error with validation errors.
func ValidationFailed(c *gin.Context, details any) {
	err := httperr.ErrValidation.WithDetails(details)
	Error(c, err)
}

// NotFound sends a 404 Not Found error.
func NotFound(c *gin.Context, resource string) {
	err := httperr.ErrVideoNotFound.WithDetails(map[string]string{
		"resource": resource,
	})
	Error(c, err)
}

// InternalError sends a 500 Internal Server Error.
func InternalError(c *gin.Context, err error) {
	ErrorFromError(c, err)
}
