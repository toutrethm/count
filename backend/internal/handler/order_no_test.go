package handler

import (
	"testing"
	"time"
)

func TestBuildOrderNoAtUsesTimestampAndSuffix(t *testing.T) {
	now := time.Date(2026, 8, 20, 18, 30, 45, 123000000, time.Local)

	got := buildOrderNoAt(now, "482")
	want := "20260820183045123-482"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
