package common

import "testing"

func TestHasRoleMatchesAllowedRole(t *testing.T) {
	if !HasRole(RoleAdmin, RoleRecruiter, RoleAdmin) {
		t.Fatal("expected admin role to be allowed")
	}
}

func TestHasRoleRejectsEmptyAndUnknownRole(t *testing.T) {
	if HasRole("", RoleAdmin) {
		t.Fatal("expected empty role to be rejected")
	}
	if HasRole(RoleJobseeker, RoleAdmin, RoleRecruiter) {
		t.Fatal("expected jobseeker role to be rejected")
	}
}

func TestSplitIDsTrimsAndDropsEmptyParts(t *testing.T) {
	got := SplitIDs(" 1, ,2,, 3 ")
	want := []string{"1", "2", "3"}

	if len(got) != len(want) {
		t.Fatalf("expected %d ids, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected id %d to be %q, got %q", i, want[i], got[i])
		}
	}
}

func TestResponseCodesStayAlignedWithHTTPStyleCodes(t *testing.T) {
	cases := map[string]int64{
		"success":      CodeSuccess,
		"param":        CodeParam,
		"unauthorized": CodeUnauthorized,
		"forbidden":    CodeForbidden,
		"not found":    CodeNotFound,
		"server":       CodeServer,
	}
	want := map[string]int64{
		"success":      200,
		"param":        400,
		"unauthorized": 401,
		"forbidden":    403,
		"not found":    404,
		"server":       500,
	}

	for name, got := range cases {
		if got != want[name] {
			t.Fatalf("expected %s code %d, got %d", name, want[name], got)
		}
	}
}
