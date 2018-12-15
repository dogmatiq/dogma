package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestStatelessAggregate_ApplyEvent_AcceptsAnyMessage(t *testing.T) {
	type message struct{} // the message interface is currently empty
	StatelessAggregate.ApplyEvent(message{})
}

func TestStatelessAggregate_ApplyEvent_PanicsOnNil(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("did not panic")
		}

		if r != "event must not be nil" {
			t.Fatal("did not panic with expected message")
		}
	}()

	StatelessAggregate.ApplyEvent(nil)
}
