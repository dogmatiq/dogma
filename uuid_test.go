package dogma

import (
	"testing"
)

func TestNormalizeUUID(t *testing.T) {
	t.Run("it returns the normalized UUID", func(t *testing.T) {
		const id = "83C4A2D9-A728-49E6-83A3-6C670B99A173"
		const want = "83c4a2d9-a728-49e6-83a3-6c670b99a173"

		got, err := normalizeUUID(id)
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Fatalf("non-normalized UUID: got %q, want %q", got, want)
		}
	})

	t.Run("it returns an error when the UUID is invalid", func(t *testing.T) {
		const valid = "b917cba9-1fa2-4513-8bf5-67acc121299f"

		cases := []string{
			valid[:len(valid)-1], // too short
			valid + "f",          // too long
		}

		// build cases for unexpected (but valid) character in each position
		for i := range len(valid) {
			id := valid
			ch := valid[i]

			if ch == '-' {
				id = id[:i] + "f" + id[i+1:] // replace hyphen with a valid hex digit
			} else {
				id = id[:i] + "-" + id[i+1:] // replace hex digit with a hyphen
			}

			cases = append(cases, id)
		}

		for _, id := range cases {
			if _, err := normalizeUUID(id); err == nil {
				t.Fatal("expected an error")
			}
		}
	})
}
