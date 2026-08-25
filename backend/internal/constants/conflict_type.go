package constants

// ConflictType is the stable vocabulary shared by geometry services and clients.
type ConflictType string

const (
	ConflictOverlap          ConflictType = "overlap"
	ConflictGap              ConflictType = "gap"
	ConflictSelfIntersection ConflictType = "self_intersection"
	ConflictDanglingEdge     ConflictType = "dangling_edge"
)

const (
	ConflictDetected           = "detected"
	ConflictConfirmed          = "confirmed"
	ConflictFalsePositive      = "false_positive"
	ConflictResolutionProposed = "resolution_proposed"
	ConflictResolved           = "resolved"
	ConflictClosed             = "closed"
)

// conflictTransitions is a directed, forward-only lifecycle: each state maps
// only to the states it may legitimately move to next. detected/confirmed/
// resolution_proposed can all be abandoned straight to closed, and a confirmed
// conflict may still be reclassified as false_positive once reviewed.
var conflictTransitions = map[string]map[string]bool{
	ConflictDetected:           {ConflictConfirmed: true, ConflictFalsePositive: true, ConflictClosed: true},
	ConflictConfirmed:          {ConflictResolutionProposed: true, ConflictFalsePositive: true, ConflictClosed: true},
	ConflictFalsePositive:      {ConflictClosed: true},
	ConflictResolutionProposed: {ConflictResolved: true, ConflictFalsePositive: true, ConflictClosed: true},
	ConflictResolved:           {ConflictClosed: true},
	ConflictClosed:             {},
}

// CanConflictTransition reports whether a conflict may move from one state to
// another. The lifecycle is directed, so the reverse direction is never legal
// merely because the forward direction is.
func CanConflictTransition(from, to string) bool {
	allowed, ok := conflictTransitions[from]
	if !ok {
		return false
	}
	return allowed[to]
}
