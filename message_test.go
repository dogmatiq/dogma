package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

func TestValidateMessage(t *testing.T) {
	err := ValidateMessage(
		fixtures.TestCommand[fixtures.TypeA]{
			Invalid: "<error>",
		},
	)
	if err == nil {
		t.Fatal("expected an error to occur")
	}

	if err.Error() != "<error>" {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestValidateMessage_nil(t *testing.T) {
	err := ValidateMessage(nil)
	if err == nil {
		t.Fatal("expected an error to occur")
	}

	if err.Error() != "message must not be nil" {
		t.Fatalf("unexpected error message: %s", err)
	}
}
