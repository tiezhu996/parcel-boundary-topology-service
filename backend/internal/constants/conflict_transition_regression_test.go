package constants

import "testing"

func TestConflictConfirmedToFalsePositive(t *testing.T) {
	if !CanConflictTransition(ConflictConfirmed, ConflictFalsePositive) {
		t.Fatal("confirmed -> false_positive should be a legal transition")
	}
}

func TestConflictProposedToClosed(t *testing.T) {
	if !CanConflictTransition(ConflictResolutionProposed, ConflictClosed) {
		t.Fatal("resolution_proposed -> closed should be a legal transition")
	}
}

func TestConflictTransitionNotSymmetric(t *testing.T) {
	if CanConflictTransition(ConflictResolutionProposed, ConflictConfirmed) {
		t.Fatal("resolution_proposed -> confirmed should be rejected (reverse transition)")
	}
}

func TestConflictDetectedToClosed(t *testing.T) {
	if !CanConflictTransition(ConflictDetected, ConflictClosed) {
		t.Fatal("detected -> closed should be a legal transition")
	}
}

func TestConflictConfirmedToClosed(t *testing.T) {
	if !CanConflictTransition(ConflictConfirmed, ConflictClosed) {
		t.Fatal("confirmed -> closed should be a legal transition")
	}
}