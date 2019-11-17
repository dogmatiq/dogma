package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
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
