package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

type executeCommandOptionBuilder struct {
	IdemKey string
}

func (b *executeCommandOptionBuilder) IdempotencyKey(key string) {
	b.IdemKey = key
}

func TestWithIdempotencyKey(t *testing.T) {
	t.Run("it returns an option with the specified idempotency key", func(t *testing.T) {
		const want = "<key>"

		b := &executeCommandOptionBuilder{}
		o := WithIdempotencyKey(want)
		o.ApplyExecuteCommandOption(b)

		if b.IdemKey != want {
			t.Fatalf("unexpected idempotency key: got %q, want %q", b.IdemKey, want)
		}
	})

	t.Run("it panics if the key is empty", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		WithIdempotencyKey("")
	})
}
