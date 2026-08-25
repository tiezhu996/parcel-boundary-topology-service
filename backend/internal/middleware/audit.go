package middleware

import (
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditContext is request-scoped metadata shared by structured request logs
// and service-level immutable audit events.
type AuditContext struct {
	RequestID string
	Method    string
	Path      string
	StartedAt time.Time
}

func AuditContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("audit_context", AuditContext{RequestID: c.GetString("request_id"), Method: c.Request.Method, Path: c.Request.URL.Path, StartedAt: time.Now().UTC()})
		c.Next()
	}
}

func AccessLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		logger.Info("http_request", "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "bytes", c.Writer.Size(), "latency_ms", time.Since(started).Milliseconds(), "client_ip", c.ClientIP(), "request_id", c.GetString("request_id"))
	}
}

// LoggerMiddleware exposes the application logger to handlers via the request
// context so that server-side errors can record their underlying cause.
func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	}
}

// LoggerFrom retrieves the application logger stored by LoggerMiddleware, or
// returns a fallback discard logger when no logger is bound (e.g. in tests).
func LoggerFrom(c *gin.Context) *slog.Logger {
	if logger, _ := c.Get("logger"); logger != nil {
		if l, ok := logger.(*slog.Logger); ok && l != nil {
			return l
		}
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
