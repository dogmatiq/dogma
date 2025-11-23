package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoCompactBehavior(t *testing.T) {
	var v NoCompactBehavior

	if err := v.Compact(t.Context(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestNoResetBehavior(t *testing.T) {
	var v NoResetBehavior

	err := v.Reset(t.Context(), nil)
	if err != ErrNotSupported {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotSupported)
	}
}

func init() {
	assertIsComparable(NoCompactBehavior{})
	assertIsComparable(NoResetBehavior{})
}
