package service

import (
	"errors"
	"testing"
	"time"

	"cadastral-boundary-topology-resolution/backend/internal/config"
	"cadastral-boundary-topology-resolution/backend/internal/dto"
	"cadastral-boundary-topology-resolution/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newAuthTestService(t *testing.T) (*AuthService, *repository.Store) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := repository.MigrateAndSeed(db); err != nil {
		t.Fatalf("migrate and seed: %v", err)
	}
	store := repository.NewStore(db)
	auth := NewAuthService(store, config.Config{JWTSecret: "auth-regression-test-secret-at-least-32-bytes", AccessTokenTTL: time.Hour})
	return auth, store
}

func TestLoginUnknownUsernameReturns401(t *testing.T) {
	auth, _ := newAuthTestService(t)
	_, err := auth.Login(dto.LoginRequest{Username: "nonexistent-user", Password: "DemoPass123!"})
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeUnauthorized {
		t.Fatalf("Login(unknown username) = %v, want CodeUnauthorized", err)
	}
}

func TestParseInactiveUserReturns401(t *testing.T) {
	auth, _ := newAuthTestService(t)
	claims := Claims{UserID: 999999, Username: "ghost", Role: "surveyor", RegisteredClaims: jwt.RegisteredClaims{Subject: "999999", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("auth-regression-test-secret-at-least-32-bytes"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	_, err = auth.Parse(token)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeUnauthorized {
		t.Fatalf("Parse(inactive user) = %v, want CodeUnauthorized", err)
	}
}
