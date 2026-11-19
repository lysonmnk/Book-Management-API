package utils

import (
	"github.com/gin-gonic/gin"
)

// APIResponse mendefinisikan format konsisten untuk response API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse mengirimkan response sukses dengan status code dan payload data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirimkan response error dengan format konsisten (misal: 400 Bad Request, 404 Not Found)
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}
