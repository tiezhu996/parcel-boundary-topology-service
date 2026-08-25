package main

import (
	"testing"
	"time"

	"cadastral-boundary-topology-resolution/backend/internal/config"
)

func TestNewLimitersUsePerEndpointLimits(t *testing.T) {
	cfg := config.Config{LoginRateLimit: 1, ImportRateLimit: 2, AnalyzeRateLimit: 3}
	login, importLimit, analyze := newLimiters(cfg)
	now := time.Now().UTC()

	if ok, _ := login.Allow("login-key", now); !ok {
		t.Fatal("first login request should be allowed")
	}
	if ok, _ := login.Allow("login-key", now); ok {
		t.Fatal("second login request should be rate-limited (login limit 1)")
	}

	if ok, _ := importLimit.Allow("import-key", now); !ok {
		t.Fatal("first import request should be allowed")
	}
	if ok, _ := importLimit.Allow("import-key", now); !ok {
		t.Fatal("second import request should be allowed (import limit 2)")
	}
	if ok, _ := importLimit.Allow("import-key", now); ok {
		t.Fatal("third import request should be rate-limited (import limit 2)")
	}

	if ok, _ := analyze.Allow("analyze-key", now); !ok {
		t.Fatal("first analyze request should be allowed")
	}
	if ok, _ := analyze.Allow("analyze-key", now); !ok {
		t.Fatal("second analyze request should be allowed")
	}
	if ok, _ := analyze.Allow("analyze-key", now); !ok {
		t.Fatal("third analyze request should be allowed (analyze limit 3)")
	}
	if ok, _ := analyze.Allow("analyze-key", now); ok {
		t.Fatal("fourth analyze request should be rate-limited (analyze limit 3)")
	}
}
