package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestWithIdempotencyKey(t *testing.T) {
	t.Run("it returns an option with the specified idempotency key", func(t *testing.T) {
		key := "test-key-123"
		option := WithIdempotencyKey(key)

		if option.IdempotencyKey() != key {
			t.Fatalf("expected idempotency key %q, got %q", key, option.IdempotencyKey())
		}
	})

	t.Run("it handles empty string", func(t *testing.T) {
		option := WithIdempotencyKey("")

		if option.IdempotencyKey() != "" {
			t.Fatalf("expected empty idempotency key, got %q", option.IdempotencyKey())
		}
	})
}

func TestExecuteCommandOption_IdempotencyKey(t *testing.T) {
	t.Run("it returns empty string for zero-value option", func(t *testing.T) {
		var option ExecuteCommandOption

		if option.IdempotencyKey() != "" {
			t.Fatalf("expected empty idempotency key, got %q", option.IdempotencyKey())
		}
	})

	t.Run("it returns the key set by WithIdempotencyKey", func(t *testing.T) {
		key := "my-unique-key"
		option := WithIdempotencyKey(key)

		if option.IdempotencyKey() != key {
			t.Fatalf("expected idempotency key %q, got %q", key, option.IdempotencyKey())
		}
	})
}
