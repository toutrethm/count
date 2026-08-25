package handler

import (
	"testing"

	"count/backend/internal/model"
)

func TestBuildOrderWorkflowPlansForGuidePillarAddsCenterHoleWhenLong(t *testing.T) {
	item := CreateOrderItemRequest{
		ComponentType: "guide_pillar",
		Quantity:      1,
		Dimensions: map[string]string{
			"length": "360",
		},
	}

	plans := buildOrderWorkflowPlans(item)

	if len(plans) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(plans))
	}
	if plans[0].ProcessCode != "blanking_center" || plans[0].StationRole != "blanking_center" {
		t.Fatalf("unexpected first step: %+v", plans[0])
	}
	if plans[1].ProcessCode != "turn_outer" || plans[1].StationRole != "turn_outer" {
		t.Fatalf("unexpected second step: %+v", plans[1])
	}
	if plans[2].ProcessCode != "turn_head" || plans[2].StationRole != "turn_head" {
		t.Fatalf("unexpected third step: %+v", plans[2])
	}
	if plans[3].ProcessCode != "center_hole" || plans[3].StationRole != "center_hole" {
		t.Fatalf("unexpected optional step: %+v", plans[3])
	}
}

func TestBuildOrderWorkflowPlansForSleeveAcceptsAutoAndManualCapabilities(t *testing.T) {
	item := CreateOrderItemRequest{
		ComponentType: "guide_bush",
		Quantity:      12,
		Dimensions: map[string]string{
			"length": "300",
		},
	}

	plans := buildOrderWorkflowPlans(item)

	if len(plans) != 1 {
		t.Fatalf("expected 1 step, got %d", len(plans))
	}
	if plans[0].ProcessCode != "turn_sleeve" {
		t.Fatalf("unexpected sleeve step: %+v", plans[0])
	}
	if plans[0].StationRole != "turn_sleeve_auto,turn_sleeve_manual" {
		t.Fatalf("expected sleeve step to accept auto and manual capability codes, got %+v", plans[0])
	}
}

func TestNextPendingOrderProcessMatchesCapabilityCodes(t *testing.T) {
	steps := []model.OrderProcess{
		{ID: 1, ProcessID: 11, Sort: 10, StationRole: "blanking_center"},
		{ID: 2, ProcessID: 12, Sort: 20, StationRole: "turn_outer"},
		{ID: 3, ProcessID: 13, Sort: 30, StationRole: "turn_head"},
	}
	scanned := map[uint]struct{}{
		1: {},
	}

	if nextPendingOrderProcess(steps, scanned, stationRoleParts("blanking_center")) != nil {
		t.Fatal("expected blanking capability to wait until turn_outer step is scanned")
	}

	got := nextPendingOrderProcess(steps, scanned, stationRoleParts("turn_outer"))
	if got == nil {
		t.Fatal("expected a pending process")
	}
	if got.ID != 2 {
		t.Fatalf("expected the next turn_outer step, got %d", got.ID)
	}

	scanned[2] = struct{}{}
	got = nextPendingOrderProcess(steps, scanned, stationRoleParts("turn_head"))
	if got == nil || got.ID != 3 {
		t.Fatalf("expected the later turn_head step, got %+v", got)
	}

	if nextPendingOrderProcess(steps, scanned, stationRoleParts("turn_sleeve_auto")) != nil {
		t.Fatal("expected nil for a capability with no pending step")
	}
}

func TestNextPendingOrderProcessSupportsMultipleCapabilities(t *testing.T) {
	steps := []model.OrderProcess{
		{ID: 1, ProcessID: 11, Sort: 10, StationRole: "blanking_center"},
		{ID: 2, ProcessID: 12, Sort: 20, StationRole: "turn_sleeve_auto,turn_sleeve_manual"},
		{ID: 3, ProcessID: 13, Sort: 30, StationRole: "turn_outer"},
	}
	scanned := map[uint]struct{}{
		1: {},
	}

	got := nextPendingOrderProcess(steps, scanned, []string{"turn_sleeve_manual", "turn_sleeve_auto"})
	if got == nil || got.ID != 2 {
		t.Fatalf("expected first allowed step for any capability, got %+v", got)
	}
}

func TestStationRolePartsSplitsCsvRoleList(t *testing.T) {
	got := stationRoleParts("turn_sleeve_auto, turn_sleeve_manual")
	if len(got) != 2 || got[0] != "turn_sleeve_auto" || got[1] != "turn_sleeve_manual" {
		t.Fatalf("unexpected station roles: %+v", got)
	}
}

func findPlan(plans []workflowStepPlan, processCode string) *workflowStepPlan {
	for i := range plans {
		if plans[i].ProcessCode == processCode {
			return &plans[i]
		}
	}
	return nil
}
