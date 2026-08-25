package service

import (
	"fmt"
	"testing"

	"cadastral-boundary-topology-resolution/backend/internal/constants"
	"cadastral-boundary-topology-resolution/backend/internal/dto"
)

func TestDetectConflictsConcurrentNeighbourParse(t *testing.T) {
	svc, _ := newCadastralTestService(t)
	surveyor := testActor(101, constants.RoleSurveyor, "002-parcel")
	base := createTestParcel(t, svc, "P-BASE", serviceTestPolygon(`[0,0],[10,0],[10,10],[0,10],[0,0]`), surveyor)
	for i := 0; i < 6; i++ {
		x0 := 10 + i*10
		x1 := x0 + 10
		pts := fmt.Sprintf(`[%d,0],[%d,0],[%d,10],[%d,10],[%d,0]`, x0, x1, x1, x0, x0)
		createTestParcel(t, svc, fmt.Sprintf("P-N%d", i), serviceTestPolygon(pts), testActor(uint(102+i), constants.RoleSurveyor, "002-neighbour"))
	}
	proposal := createTestProposal(t, svc, base, serviceTestPolygon(`[0,0],[11,0],[11,10],[0,10],[0,0]`), testActor(101, constants.RoleSurveyor, "002-proposal"))
	if _, err := svc.DetectConflicts(dto.DetectConflictRequest{ProposalID: proposal.ID, SnapToleranceM: 0.1}, "002-detect-key", testActor(101, constants.RoleGISAnalyst, "002-detect")); err != nil {
		t.Fatalf("DetectConflicts() error = %v", err)
	}
}
