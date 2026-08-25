package service

import (
	"errors"
	"testing"

	"cadastral-boundary-topology-resolution/backend/internal/dto"
)

func TestParcelNotFoundIsErrNotFound(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	_, err := svc.GetParcel(999999)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeNotFound {
		t.Fatalf("GetParcel(nonexistent) = %v, want CodeNotFound", err)
	}
}

func TestProposalNotFoundIsErrNotFound(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	_, err := svc.GetProposal(999999)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeNotFound {
		t.Fatalf("GetProposal(nonexistent) = %v, want CodeNotFound", err)
	}
}

func TestObservationNotFoundIsErrNotFound(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	_, err := svc.GetObservation(999999)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeNotFound {
		t.Fatalf("GetObservation(nonexistent) = %v, want CodeNotFound", err)
	}
}

func TestCreateParcelNewCodeSucceeds(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "003-create")
	if _, err := svc.CreateParcel(dto.CreateParcelRequest{
		ParcelCode: "P-FRESH", Name: "fresh parcel", BoundaryGeoJSON: serviceTestPolygon(`[0,0],[10,0],[10,10],[0,10],[0,0]`),
		CoordinateSystem: "EPSG:3857",
	}, actor); err != nil {
		t.Fatalf("CreateParcel(new code) error = %v, want success", err)
	}
}
