package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SuccessResponse - Return success response
func SuccessResponseJSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, SuccessResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponseJSON - Return error response
func ErrorResponseJSON(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Status:  "error",
		Message: message,
	})
}

// BadRequest - 400 Bad Request
func BadRequest(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusBadRequest, message)
}

// NotFound - 404 Not Found
func NotFound(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusNotFound, message)
}

// InternalServerError - 500 Internal Server Error
func InternalServerError(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusInternalServerError, message)
}

// Created - 201 Created
func Created(c *gin.Context, message string, data interface{}) {
	SuccessResponseJSON(c, http.StatusCreated, message, data)
}

// OK - 200 OK
func OK(c *gin.Context, message string, data interface{}) {
	SuccessResponseJSON(c, http.StatusOK, message, data)
}

// PaginatedResponse - Return paginated response
func PaginatedResponse(c *gin.Context, message string, data interface{}, pagination interface{}) {
	response := map[string]interface{}{
		"code":       http.StatusOK,
		"status":     "success",
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}
	c.JSON(http.StatusOK, response)
}

