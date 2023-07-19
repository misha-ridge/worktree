package worktree

import (
	"testing"
	"time"
)

func TestSemver(t *testing.T) {
	exp := "0.0.0-20190409104007-6aa57cbe96b8-dirty-dirtyworktree"
	ver := semver("6aa57cbe96b859c5d3d9e8ddd0a16b1e248cb7a2",
		time.Date(2019, time.April, 9, 10, 40, 7, 100, time.UTC),
		"dirtyworktree")
	if ver != exp {
		t.Errorf("Expected %s, got %s", exp, ver)
	}

	exp2 := "0.0.0-19700101000000-000000000000"
	ver2 := semver(
		"000000000000",
		time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		"")
	if ver2 != exp2 {
		t.Errorf("Expected %s, got %s", exp2, ver2)
	}

	exp3 := "0.0.0-20190409104007-6aa57cbe96b8"
	ver3 := semver("6aa57cbe96b859c5d3d9e8ddd0a16b1e248cb7a2",
		time.Date(2019, time.April, 9, 10, 40, 7, 100, time.UTC),
		"")
	if ver3 != exp3 {
		t.Errorf("Expected %s, got %s", exp3, ver3)
	}
}
