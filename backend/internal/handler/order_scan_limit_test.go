package handler

import "testing"

func TestDeriveOrderScanLimitReturnsOneForSleeveOnlyOrders(t *testing.T) {
	items := []CreateOrderItemRequest{
		{ComponentType: "guide_bush"},
		{ComponentType: "straight_sleeve"},
		{ComponentType: "middle_guide_sleeve"},
	}

	got := deriveOrderScanLimit(items)

	if got != 1 {
		t.Fatalf("expected sleeve-only order limit 1, got %d", got)
	}
}

func TestDeriveOrderScanLimitReturnsThreeWhenOrderContainsPillarComponent(t *testing.T) {
	items := []CreateOrderItemRequest{
		{ComponentType: "guide_bush"},
		{ComponentType: "guide_pillar"},
	}

	got := deriveOrderScanLimit(items)

	if got != 3 {
		t.Fatalf("expected mixed order limit 3, got %d", got)
	}
}
