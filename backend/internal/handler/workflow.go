package handler

import (
	"sort"
	"strconv"
	"strings"

	"count/backend/internal/model"
)

type workflowStepPlan struct {
	ProcessCode string
	StationRole string
	Sort        int
}

func validStationRole(role string) bool {
	return strings.TrimSpace(role) != ""
}

func normalizeStationRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	if validStationRole(role) {
		return role
	}
	return ""
}

func stationRoleParts(value string) []string {
	parts := strings.Split(value, ",")
	roles := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		role := normalizeStationRole(part)
		if role == "" {
			continue
		}
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		roles = append(roles, role)
	}
	return roles
}

func buildOrderWorkflowPlans(item CreateOrderItemRequest) []workflowStepPlan {
	componentType := strings.TrimSpace(item.ComponentType)

	switch componentType {
	case "guide_pillar", "top_pin":
		plans := []workflowStepPlan{
			{ProcessCode: "blanking_center", StationRole: "blanking_center", Sort: 10},
			{ProcessCode: "turn_outer", StationRole: "turn_outer", Sort: 20},
			{ProcessCode: "turn_head", StationRole: "turn_head", Sort: 30},
		}
		if dimensionNumber(item.Dimensions, "length") > 350 {
			plans = append(plans, workflowStepPlan{ProcessCode: "center_hole", StationRole: "center_hole", Sort: 40})
		}
		return plans
	case "b_pillar":
		return []workflowStepPlan{
			{ProcessCode: "blanking_center", StationRole: "blanking_center", Sort: 10},
			{ProcessCode: "turn_outer", StationRole: "turn_outer", Sort: 20},
			{ProcessCode: "turn_head_center", StationRole: "turn_head_center", Sort: 30},
		}
	case "pull_rod", "middle_guide_pillar":
		drillRole := "drill_tap_batch"
		drillSort := 12
		if item.Quantity < 20 {
			drillRole = "drill_tap_small"
			drillSort = 12
		}
		plans := []workflowStepPlan{
			{ProcessCode: "blanking_center", StationRole: "blanking_center", Sort: 10},
			{ProcessCode: drillRole, StationRole: drillRole, Sort: drillSort},
			{ProcessCode: "turn_outer", StationRole: "turn_outer", Sort: 20},
			{ProcessCode: "turn_head", StationRole: "turn_head", Sort: 30},
		}
		sort.SliceStable(plans, func(i, j int) bool {
			return plans[i].Sort < plans[j].Sort
		})
		return plans
	case "guide_bush", "straight_sleeve", "a_sleeve", "b_sleeve", "middle_guide_sleeve":
		return []workflowStepPlan{
			{ProcessCode: "turn_sleeve", StationRole: "turn_sleeve_auto,turn_sleeve_manual", Sort: 10},
		}
	default:
		return nil
	}
}

func dimensionNumber(dimensions map[string]string, key string) float64 {
	value := strings.TrimSpace(dimensions[key])
	if value == "" {
		return 0
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return number
}

func nextPendingOrderProcess(steps []model.OrderProcess, scanned map[uint]struct{}, stationRoles []string) *model.OrderProcess {
	allowed := make(map[string]struct{}, len(stationRoles))
	for _, role := range stationRoles {
		if normalized := normalizeStationRole(role); normalized != "" {
			allowed[normalized] = struct{}{}
		}
	}
	if len(allowed) == 0 {
		return nil
	}

	sort.SliceStable(steps, func(i, j int) bool {
		if steps[i].Sort == steps[j].Sort {
			return steps[i].ID < steps[j].ID
		}
		return steps[i].Sort < steps[j].Sort
	})

	for i := range steps {
		step := &steps[i]
		if _, ok := scanned[step.ID]; ok {
			continue
		}
		if !stepMatchesCapabilities(step.StationRole, allowed) {
			return nil
		}
		return step
	}
	return nil
}

func stepMatchesCapabilities(stepRoles string, allowed map[string]struct{}) bool {
	for _, role := range stationRoleParts(stepRoles) {
		if _, ok := allowed[role]; ok {
			return true
		}
	}
	return false
}
