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

func init() {
	assertIsComparable(NoCompactBehavior{})
}
