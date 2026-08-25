package handler

import "testing"

func TestNormalizeOrderItemModelFallsBackToPartName(t *testing.T) {
	got := normalizeOrderItemModel("", "导套")
	if got != "导套" {
		t.Fatalf("expected fallback model to part name, got %q", got)
	}
}

func TestNormalizeOrderItemModelUsesExplicitModel(t *testing.T) {
	got := normalizeOrderItemModel("M12", "导套")
	if got != "M12" {
		t.Fatalf("expected explicit model to stay, got %q", got)
	}
}
