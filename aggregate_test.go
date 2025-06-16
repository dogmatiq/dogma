package dogma_test

import (
	"testing"

	. "github.com/dogmatiq/dogma"
)

func TestNoSnapshotBehavior_TakeSnapshot_ReturnsNil(t *testing.T) {
	var v NoSnapshotBehavior

	err := v.TakeSnapshot()

	if err != nil {
		t.Fatal("unexpected error returned")
	}
}

func TestNoSnapshotBehavior_RestoreSnapshot_ReturnsNil(t *testing.T) {
	var v NoSnapshotBehavior

	err := v.RestoreSnapshot(nil)

	if err != nil {
		t.Fatal("unexpected error returned")
	}
}
