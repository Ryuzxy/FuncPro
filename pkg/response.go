package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(data any) gin.H {
	return gin.H{"status": "success", "data": data}
}

func ErrorResponse(message string) gin.H {
	return gin.H{"status": "error", "message": message}
}
