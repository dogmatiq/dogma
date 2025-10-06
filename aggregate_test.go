package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoSnapshotBehavior(t *testing.T) {
	var v NoSnapshotBehavior

	t.Run("func MarshalBinary()", func(t *testing.T) {
		t.Run("it returns ErrNotSupported", func(t *testing.T) {
			if _, err := v.MarshalBinary(); err != ErrNotSupported {
				t.Fatal(err)
			}
		})
	})

	t.Run("func UnmarshalBinary()", func(t *testing.T) {
		t.Run("it returns ErrNotSupported", func(t *testing.T) {
			if err := v.UnmarshalBinary([]byte{1, 2, 3}); err != ErrNotSupported {
				t.Fatal(err)
			}
		})
	})
}

func init() {
	assertIsComparable(NoSnapshotBehavior{})
}
