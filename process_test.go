package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

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
