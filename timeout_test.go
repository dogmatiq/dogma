package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoTimeoutHintBehavior_TimeoutHint_ReturnsZero(t *testing.T) {
	var v NoTimeoutHintBehavior

	h := v.TimeoutHint(nil)

	if h != 0 {
		t.Fatal("unexpected value returned")
	}
}
