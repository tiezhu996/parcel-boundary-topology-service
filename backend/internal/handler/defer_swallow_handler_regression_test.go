package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"cadastral-boundary-topology-resolution/backend/internal/dto"
	"cadastral-boundary-topology-resolution/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"cadastral-boundary-topology-resolution/backend/internal/model"
	"cadastral-boundary-topology-resolution/backend/internal/repository"
)

func TestTransitionObservationHandlerPreservesStatus(t *testing.T) {
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
	h := NewCadastralHandler(svc, validate)

	actor := service.Actor{ID: 101, Username: "u101", Role: "surveyor", RequestID: "007-handler"}
	parcel, err := svc.CreateParcel(dto.CreateParcelRequest{ParcelCode: "P-HD", Name: "handler", BoundaryGeoJSON: `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`, CoordinateSystem: "EPSG:3857"}, actor)
	if err != nil {
		t.Fatalf("create parcel = %v", err)
	}
	obs, err := svc.ImportObservation(dto.ImportObservationRequest{ParcelID: parcel.ID, ObservationCode: "OBS-HD", PointGeoJSON: `{"type":"Point","coordinates":[4,4]}`, ObservedAt: time.Now().UTC(), Method: "gnss", HorizontalAccuracyM: 0.1, SourceChecksum: "checksum-handler-01"}, actor)
	if err != nil {
		t.Fatalf("import = %v", err)
	}

	body, _ := json.Marshal(map[string]any{"to": "rejected", "version": obs.Version + 99})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(obs.ID))}}
	c.Set("user_id", uint(101))
	c.Set("role", "surveyor")
	c.Set("username", "u101")
	h.TransitionObservation(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("transition status = %d, want 409 (error classification preserved)", w.Code)
	}
}
