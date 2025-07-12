package dogma_test

import (
	"reflect"
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestHandlesCommand(t *testing.T) {
	t.Run("it returns a route containing the registered message type", func(t *testing.T) {
		type T struct{ Command }
		RegisterCommand[T]("83c4a2d9-a728-49e6-83a3-6c670b99a173")

		route := HandlesCommand[T]()

		got := route.Type.GoType()
		want := reflect.TypeFor[T]()

		if got != want {
			t.Fatalf("unexpected message type: got %s, want %s", got, want)
		}
	})

	t.Run("it panics if the type is not in the registry", func(t *testing.T) {
		expectPanic(
			t,
			"github.com/dogmatiq/dogma_test.T is not in the message type registry, use dogma.RegisterCommand() to add it",
			func() {
				type T struct{ Command }
				HandlesCommand[T]()
			},
		)
	})
}

func TestExecutesCommand(t *testing.T) {
	t.Run("it returns a route containing the registered message type", func(t *testing.T) {
		type T struct{ Command }
		RegisterCommand[T]("7b8cf1fd-722e-4337-bc5c-9ce4f32ab9d4")

		route := ExecutesCommand[T]()

		got := route.Type.GoType()
		want := reflect.TypeFor[T]()

		if got != want {
			t.Fatalf("unexpected message type: got %s, want %s", got, want)
		}
	})

	t.Run("it panics if the type is not in the registry", func(t *testing.T) {
		expectPanic(
			t,
			"github.com/dogmatiq/dogma_test.T is not in the message type registry, use dogma.RegisterCommand() to add it",
			func() {
				type T struct{ Command }
				ExecutesCommand[T]()
			},
		)
	})
}

func TestHandlesEvent(t *testing.T) {
	t.Run("it returns a route containing the registered message type", func(t *testing.T) {
		type T struct{ Event }
		RegisterEvent[T]("bef3014a-fca1-4cb3-90a3-ee83f5ca56c8")

		route := HandlesEvent[T]()

		got := route.Type.GoType()
		want := reflect.TypeFor[T]()

		if got != want {
			t.Fatalf("unexpected message type: got %s, want %s", got, want)
		}
	})

	t.Run("it panics if the type is not in the registry", func(t *testing.T) {
		expectPanic(
			t,
			"github.com/dogmatiq/dogma_test.T is not in the message type registry, use dogma.RegisterEvent() to add it",
			func() {
				type T struct{ Event }
				HandlesEvent[T]()
			},
		)
	})
}

func TestRecordsEvent(t *testing.T) {
	t.Run("it returns a route containing the registered message type", func(t *testing.T) {
		type T struct{ Event }
		RegisterEvent[T]("19d21601-7d10-4aaa-85b5-248cf873b3d3")

		route := RecordsEvent[T]()

		got := route.Type.GoType()
		want := reflect.TypeFor[T]()

		if got != want {
			t.Fatalf("unexpected message type: got %s, want %s", got, want)
		}
	})

	t.Run("it panics if the type is not in the registry", func(t *testing.T) {
		expectPanic(
			t,
			"github.com/dogmatiq/dogma_test.T is not in the message type registry, use dogma.RegisterEvent() to add it",
			func() {
				type T struct{ Event }
				RecordsEvent[T]()
			},
		)
	})
}

func TestSchedulesTimeout(t *testing.T) {
	t.Run("it returns a route containing the registered message type", func(t *testing.T) {
		type T struct{ Timeout }
		RegisterTimeout[T]("e11b5a92-e1ab-4a16-841a-9286b4e4d12f")

		route := SchedulesTimeout[T]()

		got := route.Type.GoType()
		want := reflect.TypeFor[T]()

		if got != want {
			t.Fatalf("unexpected message type: got %s, want %s", got, want)
		}
	})

	t.Run("it panics if the type is not in the registry", func(t *testing.T) {
		expectPanic(
			t,
			"github.com/dogmatiq/dogma_test.T is not in the message type registry, use dogma.RegisterTimeout() to add it",
			func() {
				type T struct{ Timeout }
				SchedulesTimeout[T]()
			},
		)
	})
}
