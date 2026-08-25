package middleware

import "testing"

func TestIsAllowedOriginAllowsLanHostOn5173(t *testing.T) {
	if !isAllowedOrigin("http://10.26.68.46:5173") {
		t.Fatalf("expected lan host on 5173 to be allowed")
	}
}

func TestIsAllowedOriginRejectsOtherPorts(t *testing.T) {
	if isAllowedOrigin("http://10.26.68.46:3000") {
		t.Fatalf("expected other ports to be rejected")
	}
}
