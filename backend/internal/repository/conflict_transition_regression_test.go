package repository

import (
	"testing"

	"cadastral-boundary-topology-resolution/backend/internal/constants"
	"cadastral-boundary-topology-resolution/backend/internal/dto"
	"cadastral-boundary-topology-resolution/backend/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransitionWritesNewState(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:transition-state-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.TopologyConflict{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	conflict := model.TopologyConflict{ProposalID: 1, ConflictType: constants.ConflictOverlap, ConflictState: constants.ConflictDetected, AlgorithmVersion: "v1", InputHash: "hash", Explanation: "x"}
	if err := db.Create(&conflict).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	repo := &TopologyConflictRepository{db: db}
	if err := repo.Transition(conflict.ID, constants.ConflictDetected, constants.ConflictConfirmed, nil); err != nil {
		t.Fatalf("transition: %v", err)
	}
	var got model.TopologyConflict
	if err := db.First(&got, conflict.ID).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got.ConflictState != constants.ConflictConfirmed {
		t.Fatalf("conflict state = %q, want %q", got.ConflictState, constants.ConflictConfirmed)
	}
}

func TestListConflictsFiltersState(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:list-state-filter-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.TopologyConflict{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	detected := model.TopologyConflict{ProposalID: 1, ConflictType: constants.ConflictOverlap, ConflictState: constants.ConflictDetected, AlgorithmVersion: "v1", InputHash: "a", Explanation: "x"}
	closed := model.TopologyConflict{ProposalID: 2, ConflictType: constants.ConflictOverlap, ConflictState: constants.ConflictClosed, AlgorithmVersion: "v1", InputHash: "b", Explanation: "y"}
	if err := db.Create(&detected).Error; err != nil {
		t.Fatalf("create detected: %v", err)
	}
	if err := db.Create(&closed).Error; err != nil {
		t.Fatalf("create closed: %v", err)
	}
	repo := &TopologyConflictRepository{db: db}
	items, total, err := repo.List(dto.ConflictQuery{State: constants.ConflictDetected, Page: 1, PageSize: 50})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].ConflictState != constants.ConflictDetected {
		t.Fatalf("List(state=detected) = %d items (total=%d), want exactly the detected conflict", len(items), total)
	}
}
