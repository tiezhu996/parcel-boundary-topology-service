package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRBACAllowsConfiguredRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := RBACMiddleware("surveyor")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "surveyor")
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h(c)
	if c.IsAborted() {
		t.Fatal("surveyor role should be allowed")
	}
}

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := CORSMiddleware([]string{"http://localhost:18540"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Origin", "http://localhost:18540")
	h(c)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:18540" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want the configured origin", got)
	}
}

func TestRequestIDGeneratesWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	RequestIDMiddleware()(c)
	if c.GetString("request_id") == "" {
		t.Fatal("request_id should not be empty when X-Request-ID is missing")
	}
}
