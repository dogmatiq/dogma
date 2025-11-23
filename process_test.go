package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestStatelessProcess(t *testing.T) {
	var v StatelessProcessBehavior

	root := v.New()

	if root != StatelessProcessRoot {
		t.Fatal("unexpected value returned")
	}

	t.Run("func MarshalBinary()", func(t *testing.T) {
		t.Run("it returns an empty slice", func(t *testing.T) {
			data, err := root.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if len(data) != 0 {
				t.Fatal("expected empty slice")
			}
		})
	})

	t.Run("func UnmarshalBinary()", func(t *testing.T) {
		t.Run("it returns nil if the data is empty", func(t *testing.T) {
			if err := root.UnmarshalBinary(nil); err != nil {
				t.Fatal(err)
			}

			if err := root.UnmarshalBinary([]byte{}); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("it returns an error if the data is not empty", func(t *testing.T) {
			err := root.UnmarshalBinary([]byte{1, 2, 3})
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	})
}

func TestNoTimeoutMessagesBehavior(t *testing.T) {
	var v NoTimeoutMessagesBehavior

	expectPanic(
		t,
		UnexpectedMessage,
		func() {
			v.HandleTimeout(t.Context(), nil, nil, nil)
		},
	)
}

func init() {
	assertIsComparable(StatelessProcessBehavior{})
	assertIsComparable(NoTimeoutMessagesBehavior{})
}
