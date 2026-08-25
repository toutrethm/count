package seed

import "testing"

func TestLegacyUserRoleByPhoneIncludesDrillTapForChenJiaxing(t *testing.T) {
	got := legacyUserRoleByPhone["15106017652"]
	if got != "turn_head,center_hole,drill_tap_small,drill_tap_batch,turn_sleeve_auto" {
		t.Fatalf("unexpected capability mapping: %s", got)
	}
}
