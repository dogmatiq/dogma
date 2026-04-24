package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoSnapshotBehavior(t *testing.T) {
	var v NoSnapshotBehavior

	t.Run("func MarshalBinary()", func(t *testing.T) {
		t.Run("it returns ErrNotSupported", func(t *testing.T) {
			if _, err := v.MarshalBinary(); err != ErrNotSupported {
				t.Fatal(err)
			}
		})
	})

	t.Run("func UnmarshalBinary()", func(t *testing.T) {
		t.Run("it returns ErrNotSupported", func(t *testing.T) {
			if err := v.UnmarshalBinary([]byte{1, 2, 3}); err != ErrNotSupported {
				t.Fatal(err)
			}
		})
	})
}

func TestUntypedAggregateMessageHandler(t *testing.T) {
	inner := &aggregateHandlerStub{
		routeID: "instance-001",
	}
	route := expectType[AggregateHandlerRoute](t, ViaAggregate(inner))
	adaptor := route.Handler()

	t.Run("func Configure()", func(t *testing.T) {
		t.Run("it delegates to the wrapped handler", func(t *testing.T) {
			adaptor.Configure(nil)
			if !inner.configured {
				t.Fatal("expected Configure to be called on the wrapped handler")
			}
		})
	})

	t.Run("func New()", func(t *testing.T) {
		t.Run("it returns the root from the wrapped handler", func(t *testing.T) {
			got := adaptor.New()
			expectType[aggregateRootStub](t, got)
		})
	})

	t.Run("func RouteCommandToInstance()", func(t *testing.T) {
		t.Run("it delegates to the wrapped handler", func(t *testing.T) {
			got := adaptor.RouteCommandToInstance(nil)
			if got != "instance-001" {
				t.Fatalf("unexpected instance ID: got %q, want %q", got, "instance-001")
			}
		})
	})

	t.Run("func HandleCommand()", func(t *testing.T) {
		t.Run("it narrows the root type and delegates to the wrapped handler", func(t *testing.T) {
			adaptor.HandleCommand(aggregateRootStub{}, nil, nil)
			if !inner.handleCalled {
				t.Fatal("expected HandleCommand to be called on the wrapped handler")
			}
		})

		t.Run("it panics if the root has an unexpected type", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected a panic")
				}
			}()
			adaptor.HandleCommand(nil, nil, nil)
		})
	})
}

type aggregateRootStub struct {
	NoSnapshotBehavior
}

func (aggregateRootStub) AggregateInstanceDescription() string { return "" }
func (aggregateRootStub) ApplyEvent(Event)                     {}

type aggregateHandlerStub struct {
	configured   bool
	routeID      string
	handleCalled bool
}

var _ AggregateMessageHandler[aggregateRootStub] = &aggregateHandlerStub{}

func (h *aggregateHandlerStub) Configure(AggregateConfigurer) {
	h.configured = true
}

func (h *aggregateHandlerStub) New() aggregateRootStub {
	return aggregateRootStub{}
}

func (h *aggregateHandlerStub) RouteCommandToInstance(Command) string {
	return h.routeID
}

func (h *aggregateHandlerStub) HandleCommand(
	_ aggregateRootStub,
	_ AggregateCommandScope[aggregateRootStub],
	_ Command,
) {
	h.handleCalled = true
}

func init() {
	assertIsComparable(NoSnapshotBehavior{})
}
