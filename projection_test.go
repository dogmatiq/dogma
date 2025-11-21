package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoCompactBehavior(t *testing.T) {
	var v NoCompactBehavior

	if err := v.Compact(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestResetBehavior(t *testing.T) {
	var v NoResetBehavior

	err := v.Reset(context.Background(), nil)
	if err != ErrNotSupported {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotSupported)
	}
}

func init() {
	assertIsComparable(NoCompactBehavior{})
	assertIsComparable(NoResetBehavior{})
}
