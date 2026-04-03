package dogma_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestWithEventObserver(t *testing.T) {
	t.Run("it returns an option with the registered event type", func(t *testing.T) {
		type T struct{ Event }
		RegisterEvent[*T]("b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e")

		o := WithEventObserver(func(context.Context, *T) (bool, error) {
			return true, nil
		})
		x := expectType[EventObserverOption](t, o)

		got := x.EventType().GoType()
		want := reflect.TypeFor[*T]()

		if got != want {
			t.Fatalf("unexpected event type: got %s, want %s", got, want)
		}
	})

	t.Run("it panics if the event type is not in the registry", func(t *testing.T) {
		expectPanic(
			t,
			"*github.com/dogmatiq/dogma_test.T is not in the message type registry",
			func() {
				type T struct{ Event }
				WithEventObserver(func(_ context.Context, _ *T) (bool, error) {
					return true, nil
				})
			},
		)
	})
}

func TestEventObserverOption_Observer(t *testing.T) {
	t.Run("it invokes the callback for matching events", func(t *testing.T) {
		type T struct{ Event }
		RegisterEvent[*T]("c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f")

		called := false
		wantErr := errors.New("test error")

		o := WithEventObserver(func(context.Context, *T) (bool, error) {
			called = true
			return true, wantErr
		})
		x := expectType[EventObserverOption](t, o)

		satisfied, err := x.Observer()(context.Background(), &T{})

		if !called {
			t.Fatal("expected observer to be called")
		}

		if !satisfied {
			t.Fatal("expected satisfied == true")
		}

		if err != wantErr {
			t.Fatalf("unexpected error: got %v, want %v", err, wantErr)
		}
	})

	t.Run("it ignores non-matching events", func(t *testing.T) {
		type T struct{ Event }
		RegisterEvent[*T]("d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a")

		type Other struct{ Event }
		RegisterEvent[*Other]("e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b")

		o := WithEventObserver(func(context.Context, *T) (bool, error) {
			t.Fatal("observer should not be called for non-matching event")
			return true, nil
		})
		x := expectType[EventObserverOption](t, o)

		satisfied, err := x.Observer()(context.Background(), &Other{})

		if satisfied {
			t.Fatal("expected satisfied == false for non-matching event")
		}

		if err != nil {
			t.Fatalf("unexpected error: got %v", err)
		}
	})
}
