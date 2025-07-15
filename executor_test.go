package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestWithIdempotencyKey(t *testing.T) {
	t.Run("it returns an option with the specified idempotency key", func(t *testing.T) {
		const want = "<key>"

		o := WithIdempotencyKey(want)

		x, ok := o.(IdempotencyKeyOption)
		if !ok {
			t.Fatalf("unexpected type: got %T, want %T", o, x)
		}

		if x.Key() != want {
			t.Fatalf("unexpected idempotency key: got %q, want %q", x.Key(), want)
		}
	})

	t.Run("it panics if the key is empty", func(t *testing.T) {
		expectPanic(
			t,
			`idempotency key cannot be empty`,
			func() {
				WithIdempotencyKey("")
			},
		)
	})
}
