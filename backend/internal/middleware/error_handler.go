package middleware

import (
	"errors"
	"net/http"

	"cadastral-boundary-topology-resolution/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware provides a single JSON response for errors recorded
// by middleware or future handlers that use gin.Context.Error.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}
		var appErr *service.AppError
		if errors.As(c.Errors.Last().Err, &appErr) {
			if cause := appErr.Unwrap(); cause != nil {
				LoggerFrom(c).Error("request_failed", "request_id", c.GetString("request_id"), "code", appErr.Code, "status", appErr.Status, "message", appErr.Message, "error", cause)
			}
			c.AbortWithStatusJSON(appErr.Status, gin.H{"error": gin.H{"code": appErr.Code, "message": appErr.Message, "request_id": c.GetString("request_id")}})
			return
		}
		LoggerFrom(c).Error("request_failed", "request_id", c.GetString("request_id"), "error", c.Errors.Last().Err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": service.CodeInternal, "message": "unexpected server error", "request_id": c.GetString("request_id")}})
	}
}
