package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoCompactBehavior_Compact_ReturnsNil(t *testing.T) {
	var v NoCompactBehavior

	err := v.Compact(context.Background(), nil)

	if err != nil {
		t.Fatal("unexpected error returned")
	}
}

func init() {
	assertIsComparable(NoCompactBehavior{})
}
