package handler

import "testing"

func TestCalculateWageForGuidePillarAppliesDiameterLengthAndQuantityMultiplier(t *testing.T) {
	item := CreateOrderItemRequest{
		ComponentType: "guide_pillar",
		Quantity:      12,
		Dimensions: map[string]string{
			"innerDiameter": "30",
			"length":        "220",
		},
	}

	rule, unitPrice, multiplier, amount, ok := calculateOrderItemWage(item)
	if !ok {
		t.Fatal("expected wage rule to match")
	}
	if rule.Code != "guide_pillar_d30_l151_250" {
		t.Fatalf("unexpected rule code: %s", rule.Code)
	}
	if unitPrice != 0.35 {
		t.Fatalf("unexpected unit price: %.2f", unitPrice)
	}
	if multiplier != 1.5 {
		t.Fatalf("unexpected multiplier: %.2f", multiplier)
	}
	if amount != 6.3 {
		t.Fatalf("unexpected amount: %.2f", amount)
	}
}

func TestCalculateWageForTopPinUsesSmallQuantityMultiplier(t *testing.T) {
	item := CreateOrderItemRequest{
		ComponentType: "top_pin",
		Quantity:      8,
		Dimensions: map[string]string{
			"innerDiameter": "28",
			"length":        "350",
		},
	}

	_, unitPrice, multiplier, amount, ok := calculateOrderItemWage(item)
	if !ok {
		t.Fatal("expected wage rule to match")
	}
	if unitPrice != 0.32 {
		t.Fatalf("unexpected unit price: %.2f", unitPrice)
	}
	if multiplier != 2 {
		t.Fatalf("unexpected multiplier: %.2f", multiplier)
	}
	if amount != 5.12 {
		t.Fatalf("unexpected amount: %.2f", amount)
	}
}

func TestCalculateWageForUnsupportedComponentReturnsNoMatch(t *testing.T) {
	item := CreateOrderItemRequest{
		ComponentType: "guide_bush",
		Quantity:      1,
		Dimensions: map[string]string{
			"innerDiameter": "20",
			"length":        "100",
		},
	}

	_, _, _, _, ok := calculateOrderItemWage(item)
	if ok {
		t.Fatal("expected no wage rule match")
	}
}
