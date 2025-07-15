package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

type handlerRoutesBuilder struct {
	Aggregate   ViaAggregateRoute
	Process     ViaProcessRoute
	Integration ViaIntegrationRoute
	Projection  ViaProjectionRoute
}

func (b *handlerRoutesBuilder) ViaAggregate(r ViaAggregateRoute)     { b.Aggregate = r }
func (b *handlerRoutesBuilder) ViaProcess(r ViaProcessRoute)         { b.Process = r }
func (b *handlerRoutesBuilder) ViaIntegration(r ViaIntegrationRoute) { b.Integration = r }
func (b *handlerRoutesBuilder) ViaProjection(r ViaProjectionRoute)   { b.Projection = r }

func TestViaAggregate(t *testing.T) {
	type aggregate struct{ AggregateMessageHandler }

	b := &handlerRoutesBuilder{}
	h := &aggregate{}
	r := ViaAggregate(h)
	r.ApplyHandlerRoute(b)

	if b.Aggregate.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestViaProcess(t *testing.T) {
	type process struct{ ProcessMessageHandler }

	b := &handlerRoutesBuilder{}
	h := &process{}
	r := ViaProcess(h)
	r.ApplyHandlerRoute(b)

	if b.Process.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestViaIntegration(t *testing.T) {
	type integration struct{ IntegrationMessageHandler }

	b := &handlerRoutesBuilder{}
	h := &integration{}
	r := ViaIntegration(h)
	r.ApplyHandlerRoute(b)

	if b.Integration.Handler != h {
		t.Fatal("unexpected handler")
	}
}

func TestViaProjection(t *testing.T) {
	type projection struct{ ProjectionMessageHandler }

	b := &handlerRoutesBuilder{}
	h := &projection{}
	r := ViaProjection(h)
	r.ApplyHandlerRoute(b)

	if b.Projection.Handler != h {
		t.Fatal("unexpected handler")
	}
}
