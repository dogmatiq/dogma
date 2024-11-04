package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestViaAggregate(t *testing.T) {
	type aggregate struct{ AggregateMessageHandler }

	h := &aggregate{}
	r := ViaAggregate(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestViaProcess(t *testing.T) {
	type process struct{ ProcessMessageHandler }

	h := &process{}
	r := ViaProcess(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestViaIntegration(t *testing.T) {
	type integration struct{ IntegrationMessageHandler }

	h := &integration{}
	r := ViaIntegration(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestViaProjection(t *testing.T) {
	type projection struct{ ProjectionMessageHandler }

	h := &projection{}
	r := ViaProjection(h)

	if r.Handler != h {
		t.Fatal("unexpected handler")
	}
}
