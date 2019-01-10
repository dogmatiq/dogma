package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

type testProcessRoot struct{}

func (testProcessRoot) ApplyEvent(m Message) {
	panic("not implemented")
}

func (testProcessRoot) IsEqual(r ProcessRoot) bool {
	panic("not implemented")
}

func TestStatelessProcessBehavior_New_ReturnsStatelessProcessRoot(t *testing.T) {
	var v StatelessProcessBehavior

	r := v.New()

	if r != StatelessProcessRoot {
		t.Fatal("unexpected value returned")
	}
}

func TestStatelessProcessRoot_IsEqual(t *testing.T) {
	if !StatelessProcessRoot.IsEqual(StatelessProcessRoot) {
		t.Fatal("StatelessProcessRoot is not equal to itself")
	}

	if StatelessProcessRoot.IsEqual(testProcessRoot{}) {
		t.Fatal("StatelessProcessRoot is equal to a different process root")
	}
}

func TestStatelessProcessRoot_IsEqual_PanicsOnNil(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("did not panic")
		}

		if r != "process root must not be nil" {
			t.Fatal("did not panic with expected message")
		}
	}()

	StatelessProcessRoot.IsEqual(nil)
}

func TestNoTimeoutBehavior_HandleTimeout_Panics(t *testing.T) {
	var v NoTimeoutBehavior
	ctx := context.Background()

	defer func() {
		r := recover()

		if r != UnexpectedMessage {
			t.Fatal("expected panic did not occur")
		}
	}()

	v.HandleTimeout(ctx, nil, nil)
}
