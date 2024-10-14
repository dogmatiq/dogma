package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestRegisterAggregate(t *testing.T) {
	type aggregate struct{ AggregateMessageHandler }

	h := &aggregate{}
	r := RegisterAggregate(h)

	if r, ok := r.(AggregateRegistration); ok {
		if r.Handler != h {
			t.Fatal("unexpected handler")
		}
	} else {
		t.Fatal("unexpected type")
	}
}

func TestRegisterProcess(t *testing.T) {
	type process struct{ ProcessMessageHandler }

	h := &process{}
	r := RegisterProcess(h)

	if r, ok := r.(ProcessRegistration); ok {
		if r.Handler != h {
			t.Fatal("unexpected handler")
		}
	} else {
		t.Fatal("unexpected type")
	}
}

func TestRegisterIntegration(t *testing.T) {
	type integration struct{ IntegrationMessageHandler }

	h := &integration{}
	r := RegisterIntegration(h)

	if r, ok := r.(IntegrationRegistration); ok {
		if r.Handler != h {
			t.Fatal("unexpected handler")
		}
	} else {
		t.Fatal("unexpected type")
	}
}

func TestRegisterProjection(t *testing.T) {
	type projection struct{ ProjectionMessageHandler }

	h := &projection{}
	r := RegisterProjection(h)

	if r, ok := r.(ProjectionRegistration); ok {
		if r.Handler != h {
			t.Fatal("unexpected handler")
		}
	} else {
		t.Fatal("unexpected type")
	}
}
