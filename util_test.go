package dogma_test

import (
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
func expectPanic(t *testing.T, message string, fn func()) {
	t.Helper()

	defer func() {
		switch got := recover(); got {
		case nil:
			t.Fatal("expected function to panic")
		case message:
			// ok
		default:
			t.Fatalf(
				"unexpected panic message:\n  got %q (%T),\n want %q",
				got,
				got,
				message,
			)
		}
	}()

	fn()
}
