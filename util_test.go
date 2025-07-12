package dogma_test

import (
	"testing"
)

// messageWithPointerRecievers is an implementation of [Message] that uses
// pointer receivers, used to test reflection code that has pointer-specific
// logic.
//
// S is the validation scope type accepted by the Validate() method, which lets
// the implementation match any of the [Command], [Event], or [Timeout]
// interfaces as required.
type messageWithPointerRecievers[S any] struct{}

func (*messageWithPointerRecievers[S]) MessageDescription() string { panic("not implemented") }
func (*messageWithPointerRecievers[S]) Validate(S) error           { panic("not implemented") }

// expectPanic is a test helper that asserts that fn panics with a specific
// message.
func expectPanic(t *testing.T, message string, fn func()) {
	t.Helper()

	defer func() {
		switch got := recover(); got {
		case nil:
			t.Fatal("expected panic")
		case message:
			// ok
		default:
			t.Fatalf("unexpected panic message: got %q, want %q", got, message)
		}
	}()

	fn()
}
