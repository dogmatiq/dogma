package dogma_test

import (
	"fmt"
	"testing"
)

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
