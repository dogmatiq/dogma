package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestStatelessProcess_New_ReturnsStatelessProcessRoot(t *testing.T) {
	var v StatelessProcess

	r := v.New()

	if r != StatelessProcessRoot {
		t.Fatal("unexpected value returned")
	}
}

func TestNoTimeouts_HandleTimeout_Panics(t *testing.T) {
	var v NoTimeouts
	ctx := context.Background()

	defer func() {
		r := recover()

		if r != UnexpectedMessage {
			t.Fatal("expected panic did not occur")
		}
	}()

	v.HandleTimeout(ctx, nil, nil)
}
