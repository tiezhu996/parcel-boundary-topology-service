package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"cadastral-boundary-topology-resolution/backend/internal/model"
	"cadastral-boundary-topology-resolution/backend/internal/repository"
	"cadastral-boundary-topology-resolution/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newHandler(t *testing.T) *CadastralHandler {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.AuditLog{}, &model.LandParcel{}, &model.SurveyObservation{}, &model.BoundaryProposal{}, &model.TopologyConflict{}, &model.TopologyDetectionRun{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	store := repository.NewStore(db)
	svc := service.NewCadastralService(store)
	validate := validator.New(validator.WithRequiredStructEnabled())
	return NewCadastralHandler(svc, validate)
}

func TestGetConflictNotFoundStatus(t *testing.T) {
	h := newHandler(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "999999"}}
	h.GetConflict(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("GetConflict status = %d, want 404", w.Code)
	}
}

func TestTransitionConflictNotFoundStatus(t *testing.T) {
	h := newHandler(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"to":"confirmed"}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "999999"}}
	c.Set("user_id", uint(101))
	c.Set("role", "reviewer")
	c.Set("username", "u101")
	h.TransitionConflict(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("TransitionConflict status = %d, want 404", w.Code)
	}
}

func TestGetParcelNotFoundStatus(t *testing.T) {
	h := newHandler(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "999999"}}
	h.GetParcel(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("GetParcel status = %d, want 404", w.Code)
	}
}

func TestUpdateParcelNotFoundStatus(t *testing.T) {
	h := newHandler(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader([]byte(`{"name":"xy"}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "999999"}}
	c.Set("user_id", uint(101))
	c.Set("role", "surveyor")
	c.Set("username", "u101")
	h.UpdateParcel(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("UpdateParcel status = %d, want 404", w.Code)
	}
}

func TestDetectConflictsNotFoundStatus(t *testing.T) {
	h := newHandler(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"proposal_id":999999,"snap_tolerance_m":0.1}`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Idempotency-Key", "detect-key-010")
	c.Set("user_id", uint(101))
	c.Set("role", "gis_analyst")
	c.Set("username", "u101")
	h.DetectConflicts(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("DetectConflicts status = %d, want 404", w.Code)
	}
}
