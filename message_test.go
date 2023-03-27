package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

func TestDescribeMessage(t *testing.T) {
	expect := "command(int:123, valid)"
	actual := DescribeMessage(
		fixtures.Command[int]{
			Content: 123,
		},
	)
	if actual != expect {
		t.Fatalf(
			"unexpected message description: want %q, got %q",
			expect,
			actual,
		)
	}
}

func TestValidateMessage_validatable(t *testing.T) {
	expect := "<error>"
	err := ValidateMessage(
		fixtures.Command[int]{
			Invalid: expect,
		},
	)
	if err == nil {
		t.Fatal("expected an error to occur")
	}
	if err.Error() != expect {
		t.Fatalf(
			"unexpected error: want %q, got %q",
			expect,
			err,
		)
	}
}

func TestValidateMessage_nil(t *testing.T) {
	expect := "message must not be nil"
	err := ValidateMessage(nil)
	if err == nil {
		t.Fatal("expected an error to occur")
	}
	if err.Error() != expect {
		t.Fatalf(
			"unexpected error message: want %q, got %q",
			expect,
			err,
		)
	}
}

func TestValidateMessage_default(t *testing.T) {
	type desc interface {
		MessageDescription() string
	}

	err := ValidateMessage(
		struct{ desc }{},
	)
	if err != nil {
		t.Fatal("unexpected error")
	}
}
