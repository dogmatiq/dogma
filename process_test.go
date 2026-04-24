package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestStatelessProcess(t *testing.T) {
	var v StatelessProcessBehavior

	root := v.New()

	if root != (StatelessProcessRoot{}) {
		t.Fatal("unexpected value returned")
	}

	t.Run("func ProcessInstanceDescription()", func(t *testing.T) {
		t.Run("it returns an empty string", func(t *testing.T) {
			if got := root.ProcessInstanceDescription(false); got != "" {
				t.Fatalf("unexpected description: %q", got)
			}
		})

		t.Run("it returns an empty string when the instance has ended", func(t *testing.T) {
			if got := root.ProcessInstanceDescription(true); got != "" {
				t.Fatalf("unexpected description: %q", got)
			}
		})
	})

	t.Run("func MarshalBinary()", func(t *testing.T) {
		t.Run("it returns an empty slice", func(t *testing.T) {
			data, err := root.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if len(data) != 0 {
				t.Fatal("expected empty slice")
			}
		})
	})

	t.Run("func UnmarshalBinary()", func(t *testing.T) {
		t.Run("it returns nil if the data is empty", func(t *testing.T) {
			if err := root.UnmarshalBinary(nil); err != nil {
				t.Fatal(err)
			}

			if err := root.UnmarshalBinary([]byte{}); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("it returns an error if the data is not empty", func(t *testing.T) {
			err := root.UnmarshalBinary([]byte{1, 2, 3})
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	})
}

func TestNoTimeoutMessagesBehavior(t *testing.T) {
	var v NoTimeoutMessagesBehavior[ProcessRoot]

	expectPanic(
		t,
		UnexpectedMessage,
		func() {
			v.HandleTimeout(t.Context(), nil, nil, nil)
		},
	)
}

func TestUntypedProcessMessageHandler(t *testing.T) {
	inner := &processHandlerStub{
		routeID: "instance-001",
		routeOK: true,
	}
	route := expectType[ProcessHandlerRoute](t, ViaProcess(inner))
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
			expectType[processRootStub](t, got)
		})
	})

	t.Run("func RouteEventToInstance()", func(t *testing.T) {
		t.Run("it delegates to the wrapped handler", func(t *testing.T) {
			id, ok, err := adaptor.RouteEventToInstance(t.Context(), nil)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Fatal("expected ok to be true")
			}
			if id != "instance-001" {
				t.Fatalf("unexpected instance ID: got %q, want %q", id, "instance-001")
			}
		})
	})

	t.Run("func HandleEvent()", func(t *testing.T) {
		t.Run("it narrows the root type and delegates to the wrapped handler", func(t *testing.T) {
			err := adaptor.HandleEvent(t.Context(), processRootStub{}, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !inner.handleEventCalled {
				t.Fatal("expected HandleEvent to be called on the wrapped handler")
			}
		})

		t.Run("it panics if the root has an unexpected type", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected a panic")
				}
			}()
			adaptor.HandleEvent(t.Context(), nil, nil, nil)
		})
	})

	t.Run("func HandleTimeout()", func(t *testing.T) {
		t.Run("it narrows the root type and delegates to the wrapped handler", func(t *testing.T) {
			err := adaptor.HandleTimeout(t.Context(), processRootStub{}, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !inner.handleTimeoutCalled {
				t.Fatal("expected HandleTimeout to be called on the wrapped handler")
			}
		})

		t.Run("it panics if the root has an unexpected type", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected a panic")
				}
			}()
			adaptor.HandleTimeout(t.Context(), nil, nil, nil)
		})
	})
}

type processRootStub struct{}

func (processRootStub) ProcessInstanceDescription(bool) string { return "" }
func (processRootStub) MarshalBinary() ([]byte, error)         { return nil, nil }
func (processRootStub) UnmarshalBinary([]byte) error           { return nil }

type processHandlerStub struct {
	configured          bool
	routeID             string
	routeOK             bool
	handleEventCalled   bool
	handleTimeoutCalled bool
}

var _ ProcessMessageHandler[processRootStub] = &processHandlerStub{}

func (h *processHandlerStub) Configure(ProcessConfigurer) {
	h.configured = true
}

func (h *processHandlerStub) New() processRootStub {
	return processRootStub{}
}

func (h *processHandlerStub) RouteEventToInstance(
	_ context.Context,
	_ Event,
) (string, bool, error) {
	return h.routeID, h.routeOK, nil
}

func (h *processHandlerStub) HandleEvent(
	_ context.Context,
	_ processRootStub,
	_ ProcessEventScope[processRootStub],
	_ Event,
) error {
	h.handleEventCalled = true
	return nil
}

func (h *processHandlerStub) HandleTimeout(
	_ context.Context,
	_ processRootStub,
	_ ProcessTimeoutScope[processRootStub],
	_ Timeout,
) error {
	h.handleTimeoutCalled = true
	return nil
}

func init() {
	assertIsComparable(StatelessProcessBehavior{})
	assertIsComparable(NoTimeoutMessagesBehavior[ProcessRoot]{})
	assertIsComparable(StatelessProcessRoot{})
}
