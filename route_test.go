package dogma_test

import (
	"reflect"
	"testing"

	. "github.com/dogmatiq/dogma"
)

type nonPointerReceivers[S any] struct{}
type pointerReceivers[S any] struct{}

func (nonPointerReceivers[S]) MessageDescription() string { panic("not implemented") }
func (nonPointerReceivers[S]) Validate(S) error           { panic("not implemented") }
func (*pointerReceivers[S]) MessageDescription() string   { panic("not implemented") }
func (*pointerReceivers[S]) Validate(S) error             { panic("not implemented") }

func TestHandlesCommand(t *testing.T) {
	type (
		N = nonPointerReceivers[CommandValidationScope]
		P = *pointerReceivers[CommandValidationScope]
		X = *nonPointerReceivers[CommandValidationScope]
	)

	t.Run("it returns a route with the correct reflection type", func(t *testing.T) {
		if HandlesCommand[N]().Type != reflect.TypeFor[N]() {
			t.Fatal("unexpected message type")
		}

		if HandlesCommand[P]().Type != reflect.TypeFor[P]() {
			t.Fatal("unexpected message type")
		}
	})

	t.Run("it panics if the type is a pointer to an implementation that uses non-pointer receivers", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()
		HandlesCommand[X]()
	})
}

func TestRecordsEvent(t *testing.T) {
	type (
		N = nonPointerReceivers[EventValidationScope]
		P = *pointerReceivers[EventValidationScope]
		X = *nonPointerReceivers[EventValidationScope]
	)

	t.Run("it returns a route with the correct reflection type", func(t *testing.T) {
		if RecordsEvent[N]().Type != reflect.TypeFor[N]() {
			t.Fatal("unexpected message type")
		}

		if RecordsEvent[P]().Type != reflect.TypeFor[P]() {
			t.Fatal("unexpected message type")
		}
	})

	t.Run("it panics if the type is a pointer to an implementation that uses non-pointer receivers", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()
		RecordsEvent[X]()
	})
}

func TestHandlesEvent(t *testing.T) {
	type (
		N = nonPointerReceivers[EventValidationScope]
		P = *pointerReceivers[EventValidationScope]
		X = *nonPointerReceivers[EventValidationScope]
	)

	t.Run("it returns a route with the correct reflection type", func(t *testing.T) {
		if HandlesEvent[N]().Type != reflect.TypeFor[N]() {
			t.Fatal("unexpected message type")
		}

		if HandlesEvent[P]().Type != reflect.TypeFor[P]() {
			t.Fatal("unexpected message type")
		}
	})

	t.Run("it panics if the type is a pointer to an implementation that uses non-pointer receivers", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()
		HandlesEvent[X]()
	})
}

func TestExecutesCommand(t *testing.T) {
	type (
		N = nonPointerReceivers[CommandValidationScope]
		P = *pointerReceivers[CommandValidationScope]
		X = *nonPointerReceivers[CommandValidationScope]
	)

	t.Run("it returns a route with the correct reflection type", func(t *testing.T) {
		if ExecutesCommand[N]().Type != reflect.TypeFor[N]() {
			t.Fatal("unexpected message type")
		}

		if ExecutesCommand[P]().Type != reflect.TypeFor[P]() {
			t.Fatal("unexpected message type")
		}
	})

	t.Run("it panics if the type is a pointer to an implementation that uses non-pointer receivers", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()
		ExecutesCommand[X]()
	})
}

func TestSchedulesTimeout(t *testing.T) {
	type (
		N = nonPointerReceivers[TimeoutValidationScope]
		P = *pointerReceivers[TimeoutValidationScope]
		X = *nonPointerReceivers[TimeoutValidationScope]
	)

	t.Run("it returns a route with the correct reflection type", func(t *testing.T) {
		if SchedulesTimeout[N]().Type != reflect.TypeFor[N]() {
			t.Fatal("unexpected message type")
		}

		if SchedulesTimeout[P]().Type != reflect.TypeFor[P]() {
			t.Fatal("unexpected message type")
		}
	})

	t.Run("it panics if the type is a pointer to an implementation that uses non-pointer receivers", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()
		SchedulesTimeout[X]()
	})
}
