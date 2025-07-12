package dogma_test

import (
	"testing"
)

type messageWithPointerRecievers[S any] struct{}

func (*messageWithPointerRecievers[S]) MessageDescription() string { panic("not implemented") }
func (*messageWithPointerRecievers[S]) Validate(S) error           { panic("not implemented") }

func expectPanic(t *testing.T, want string, fn func()) {
	defer func() {
		switch got := recover(); got {
		case nil:
			t.Fatal("expected panic")
		case want:
			// ok
		default:
			t.Fatalf("unexpected panic message: got %q, want %q", got, want)
		}
	}()

	fn()
}
