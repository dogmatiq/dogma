package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestWithIdempotencyKey(t *testing.T) {
	t.Run("it returns an option with the specified idempotency key", func(t *testing.T) {
		key := "test-key-123"
		option := WithIdempotencyKey(key)

		if IdempotencyKey(option) != key {
			t.Fatalf("expected idempotency key %q, got %q", key, IdempotencyKey(option))
		}
	})

	t.Run("it handles empty string", func(t *testing.T) {
		option := WithIdempotencyKey("")

		if IdempotencyKey(option) != "" {
			t.Fatalf("expected empty idempotency key, got %q", IdempotencyKey(option))
		}
	})
}

func TestIdempotencyKey(t *testing.T) {
	t.Run("it returns empty string for no options", func(t *testing.T) {
		if IdempotencyKey() != "" {
			t.Fatalf("expected empty idempotency key, got %q", IdempotencyKey())
		}
	})

	t.Run("it returns the key from WithIdempotencyKey option", func(t *testing.T) {
		key := "my-unique-key"
		option := WithIdempotencyKey(key)

		if IdempotencyKey(option) != key {
			t.Fatalf("expected idempotency key %q, got %q", key, IdempotencyKey(option))
		}
	})

	t.Run("it returns the key from multiple options", func(t *testing.T) {
		key := "my-unique-key"
		option1 := WithIdempotencyKey(key)
		
		if IdempotencyKey(option1) != key {
			t.Fatalf("expected idempotency key %q, got %q", key, IdempotencyKey(option1))
		}
	})
}
