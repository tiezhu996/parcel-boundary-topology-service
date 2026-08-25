package service

import (
	"errors"
	"testing"

	"cadastral-boundary-topology-resolution/backend/internal/dto"
)

const selfIntersecting = `{"type":"Polygon","coordinates":[[[0,0],[10,10],[0,10],[10,0],[0,0]]]}`

func TestCreateProposalInvalidGeometryIs422(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "004-prop")
	parcel := createTestParcel(t, svc, "P-004", serviceTestPolygon(`[0,0],[10,0],[10,10],[0,10],[0,0]`), actor)
	_, err := svc.CreateProposal(dto.CreateProposalRequest{
		ParcelID: parcel.ID, BaseVersion: parcel.BoundaryVersion,
		ProposedGeoJSON: selfIntersecting, SnapToleranceM: 0.1,
	}, actor)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeInvalidInput {
		t.Fatalf("CreateProposal(invalid geometry) = %v, want CodeInvalidInput", err)
	}
}

func TestCreateParcelInvalidGeometryIs422(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "004-parcel")
	_, err := svc.CreateParcel(dto.CreateParcelRequest{
		ParcelCode: "P-BAD", Name: "bad", BoundaryGeoJSON: selfIntersecting, CoordinateSystem: "EPSG:3857",
	}, actor)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeInvalidInput {
		t.Fatalf("CreateParcel(invalid geometry) = %v, want CodeInvalidInput", err)
	}
}

func TestUpdateParcelInvalidGeometryIs422(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "004-update")
	parcel := createTestParcel(t, svc, "P-UPD", serviceTestPolygon(`[0,0],[10,0],[10,10],[0,10],[0,0]`), actor)
	bad := selfIntersecting
	_, err := svc.UpdateParcel(parcel.ID, dto.UpdateParcelRequest{BoundaryGeoJSON: &bad}, actor)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeInvalidInput {
		t.Fatalf("UpdateParcel(invalid geometry) = %v, want CodeInvalidInput", err)
	}
}

func TestCreateParcelInvalidCRSIs422(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	actor := testActor(101, "surveyor", "004-crs")
	_, err := svc.CreateParcel(dto.CreateParcelRequest{
		ParcelCode: "P-CRS", Name: "crs", BoundaryGeoJSON: serviceTestPolygon(`[0,0],[10,0],[10,10],[0,10],[0,0]`),
		CoordinateSystem: "EPSG:4326",
	}, actor)
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Code != CodeInvalidInput {
		t.Fatalf("CreateParcel(invalid CRS) = %v, want CodeInvalidInput", err)
	}
}
