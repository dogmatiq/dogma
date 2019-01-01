package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestStatelessAggregate_New_ReturnsStatelessAggregateRoot(t *testing.T) {
	var v StatelessAggregate

	r := v.New()

	if r != StatelessAggregateRoot {
		t.Fatal("unexpected value returned")
	}
}

func TestStatelessAggregateRoot_ApplyEvent_AcceptsAnyMessage(t *testing.T) {
	type message struct{} // the message interface is currently empty
	StatelessAggregateRoot.ApplyEvent(message{})
}

func TestStatelessAggregateRoot_ApplyEvent_PanicsOnNil(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("did not panic")
		}

		if r != "event must not be nil" {
			t.Fatal("did not panic with expected message")
		}
	}()

	StatelessAggregateRoot.ApplyEvent(nil)
}
