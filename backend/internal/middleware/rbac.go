package middleware

import (
	"net/http"

	"cadastral-boundary-topology-resolution/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(roles ...string) gin.HandlerFunc {
	var allowed map[string]bool
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		roleString, isString := role.(string)
		if !ok || !isString || !allowed[roleString] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": service.CodeForbidden, "message": "role is not permitted for this operation", "request_id": c.GetString("request_id")}})
			return
		}
		c.Next()
	}
}
