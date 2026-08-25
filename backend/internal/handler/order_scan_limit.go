package handler

import (
	"strings"

	"count/backend/internal/model"
)

func deriveOrderScanLimit(items []CreateOrderItemRequest) int {
	limit := 1
	for _, item := range items {
		if componentScanLimit(item.ComponentType) > limit {
			limit = 3
		}
	}
	return limit
}

func orderScanLimitForOrder(order model.Order) int {
	if order.ScanLimit > 0 {
		return order.ScanLimit
	}

	limit := 1
	for _, item := range order.Items {
		if componentScanLimit(item.ComponentType) > limit {
			limit = 3
		}
	}
	return limit
}

func componentScanLimit(componentType string) int {
	switch strings.TrimSpace(componentType) {
	case "guide_bush", "straight_sleeve", "a_sleeve", "b_sleeve", "middle_guide_sleeve":
		return 1
	default:
		return 3
	}
}
