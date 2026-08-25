package service

import (
	"errors"
	"testing"
	"time"

	"cadastral-boundary-topology-resolution/backend/internal/dto"
)

func TestImportObservationPreservesError(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "007-import")
	_, err := svc.ImportObservation(dto.ImportObservationRequest{ParcelID: 999999, ObservationCode: "OBS-X", PointGeoJSON: `{"type":"Point","coordinates":[1,1]}`, ObservedAt: time.Now().UTC(), Method: "gnss", HorizontalAccuracyM: 0.1, SourceChecksum: "checksum-12345678"}, actor)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Err == nil {
		t.Fatalf("ImportObservation error should preserve the underlying cause, got %v", err)
	}
}

func TestGetObservationPreservesError(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	_, err := svc.GetObservation(999999)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Err == nil {
		t.Fatalf("GetObservation error should preserve the underlying cause, got %v", err)
	}
}

func TestTransitionObservationPreservesError(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "007-trans")
	_, err := svc.TransitionObservation(999999, dto.ObservationTransitionRequest{To: "rejected", Version: 1}, actor)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Err == nil {
		t.Fatalf("TransitionObservation error should preserve the underlying cause, got %v", err)
	}
}
