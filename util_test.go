package dogma_test

import (
	"fmt"
	"testing"
)

// messageWithPointerReceivers is an implementation of [Message] that uses
// pointer receivers.
//
// S is the validation scope type accepted by the Validate() method, which lets
// the implementation match any of the [Command], [Event], or [Timeout]
// interfaces as required.
type messageWithPointerReceivers[S any] struct{}

func (*messageWithPointerReceivers[S]) MessageDescription() string { panic("not implemented") }
func (*messageWithPointerReceivers[S]) Validate(S) error           { panic("not implemented") }

// expectPanic is a test helper that asserts that fn panics with a specific
// message.
func expectPanic[T comparable](t *testing.T, message T, fn func()) {
	t.Helper()

	var got any

	func() {
		defer func() { got = recover() }()
		fn()
	}()

	switch got {
	case nil:
		t.Fatal("expected function to panic")
	case message:
		// ok
	default:
		t.Fatalf(
			"unexpected panic message:\n  got %q (%T),\n want %q",
			got,
			got,
			fmt.Sprintf("%v", message),
		)
	}
}

func expectType[T any](t *testing.T, v any) T {
	t.Helper()

	got, ok := v.(T)

	if !ok {
		var want T
		t.Fatalf("unexpected type: got %T, want %T", got, want)
	}

	return got
}

func assertIsComparable[T comparable](T) {}
