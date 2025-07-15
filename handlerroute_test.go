package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestViaAggregate(t *testing.T) {
	t.Run("it returns a route with the specified handler", func(t *testing.T) {
		type aggregate struct{ AggregateMessageHandler }

		h := &aggregate{}
		r := ViaAggregate(h)
		x := expectType[AggregateHandlerRoute](t, r)

		if x.Handler() != h {
			t.Fatal("unexpected handler")
		}
	})

	t.Run("it panics if the handler is nil", func(t *testing.T) {
		expectPanic(
			t,
			`handler cannot be nil`,
			func() {
				ViaAggregate(nil)
			},
		)
	})
}

func TestViaProcess(t *testing.T) {
	t.Run("it returns a route with the specified handler", func(t *testing.T) {
		type process struct{ ProcessMessageHandler }

		h := &process{}
		r := ViaProcess(h)
		x := expectType[ProcessHandlerRoute](t, r)

		if x.Handler() != h {
			t.Fatal("unexpected handler")
		}
	})

	t.Run("it panics if the handler is nil", func(t *testing.T) {
		expectPanic(
			t,
			`handler cannot be nil`,
			func() {
				ViaProcess(nil)
			},
		)
	})
}

func TestViaIntegration(t *testing.T) {
	t.Run("it returns a route with the specified handler", func(t *testing.T) {
		type integration struct{ IntegrationMessageHandler }

		h := &integration{}
		r := ViaIntegration(h)
		x := expectType[IntegrationHandlerRoute](t, r)

		if x.Handler() != h {
			t.Fatal("unexpected handler")
		}
	})

	t.Run("it panics if the handler is nil", func(t *testing.T) {
		expectPanic(
			t,
			`handler cannot be nil`,
			func() {
				ViaIntegration(nil)
			},
		)
	})
}

func TestViaProjection(t *testing.T) {
	t.Run("it returns a route with the specified handler", func(t *testing.T) {
		type projection struct{ ProjectionMessageHandler }

		h := &projection{}
		r := ViaProjection(h)
		x := expectType[ProjectionHandlerRoute](t, r)

		if x.Handler() != h {
			t.Fatal("unexpected handler")
		}
	})

	t.Run("it panics if the handler is nil", func(t *testing.T) {
		expectPanic(
			t,
			`handler cannot be nil`,
			func() {
				ViaProjection(nil)
			},
		)
	})
}
