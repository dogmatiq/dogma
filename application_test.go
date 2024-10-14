package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestRegisterAggregate(t *testing.T) {
	type aggregate struct{ AggregateMessageHandler }

	h := &aggregate{}
	r := RegisterAggregate(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestRegisterProcess(t *testing.T) {
	type process struct{ ProcessMessageHandler }

	h := &process{}
	r := RegisterProcess(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestRegisterIntegration(t *testing.T) {
	type integration struct{ IntegrationMessageHandler }

	h := &integration{}
	r := RegisterIntegration(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestRegisterProjection(t *testing.T) {
	type projection struct{ ProjectionMessageHandler }

	h := &projection{}
	r := RegisterProjection(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}
