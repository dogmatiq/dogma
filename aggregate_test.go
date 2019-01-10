package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

type testAggregateRoot struct{}

func (testAggregateRoot) ApplyEvent(m Message) {
	panic("not implemented")
}

func (testAggregateRoot) IsEqual(r AggregateRoot) bool {
	panic("not implemented")
}

func TestStatelessAggregateBehavior_New_ReturnsStatelessAggregateRoot(t *testing.T) {
	var v StatelessAggregateBehavior

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

func TestStatelessAggregateRoot_IsEqual(t *testing.T) {
	if !StatelessAggregateRoot.IsEqual(StatelessAggregateRoot) {
		t.Fatal("StatelessAggregateRoot is not equal to itself")
	}

	if StatelessAggregateRoot.IsEqual(testAggregateRoot{}) {
		t.Fatal("StatelessAggregateRoot is equal to a different aggregate root")
	}
}

func TestStatelessAggregateRoot_IsEqual_PanicsOnNil(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("did not panic")
		}

		if r != "aggregate root must not be nil" {
			t.Fatal("did not panic with expected message")
		}
	}()

	StatelessAggregateRoot.IsEqual(nil)
}
