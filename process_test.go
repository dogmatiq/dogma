package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestStatelessProcessBehavior_New_ReturnsStatelessProcessRoot(t *testing.T) {
	var v StatelessProcessBehavior

	r := v.New()

	if r != StatelessProcessRoot {
		t.Fatal("unexpected value returned")
	}
}

func TestNoTimeoutMessagesBehavior_HandleTimeout_Panics(t *testing.T) {
	var v NoTimeoutMessagesBehavior
	ctx := context.Background()

	expectPanic(
		t,
		UnexpectedMessage,
		func() {
			v.HandleTimeout(ctx, nil, nil, nil)
		},
	)
}

func init() {
	assertIsComparable(StatelessProcessBehavior{})
	assertIsComparable(NoTimeoutMessagesBehavior{})
}
