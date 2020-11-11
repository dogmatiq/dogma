package dogma_test

import (
	"errors"
	"testing"

	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

type describable struct{}

func (describable) MessageDescription() string {
	return "<description>"
}

func (describable) String() string {
	panic("unexpected call")
}

type stringer struct{}

func (stringer) String() string {
	return "<string>"
}

type indescribable struct {
	Value int
}

func TestDescribeMessage_describable(t *testing.T) {
	d := DescribeMessage(describable{})
	if d != "<description>" {
		t.Fatal("unexpected message description")
	}
}

func TestDescribeMessage_stringer(t *testing.T) {
	d := DescribeMessage(stringer{})
	if d != "<string>" {
		t.Fatal("unexpected message description")
	}
}

func TestDescribeMessage_default(t *testing.T) {
	d := DescribeMessage(indescribable{100})
	if d != "{100}" {
		t.Fatal("unexpected message description")
	}
}

func TestValidateMessage_validatable(t *testing.T) {
	expect := errors.New("<error>")

	err := ValidateMessage(fixtures.MessageA{
		Value: expect,
	})
	if err == nil {
		t.Fatal("expected an error to occur")
	}

	if err != expect {
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

func TestValidateMessage_default(t *testing.T) {
	err := ValidateMessage(struct{}{})
	if err != nil {
		t.Fatal("unexpected error")
	}
}
