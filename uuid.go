package dogma

import (
	"fmt"
	"uuid"
)

// parseUUID parses s as a canonical RFC 9562 UUID, such as
// "65f9620a-65c1-434e-8292-60cd7938c4de", in any mix of upper and lower case.
//
// It rejects the other forms that [uuid.Parse] accepts, and the nil UUID.
func parseUUID(s string) (uuid.UUID, error) {
	// The canonical form is the only form of this length that [uuid.Parse]
	// accepts, so checking the length first rejects the other forms.
	if len(s) != 36 {
		return uuid.UUID{}, fmt.Errorf("%q is not a canonical RFC 9562 UUID: expected 36 characters", s)
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%q is not a canonical RFC 9562 UUID: expected hyphens and hexadecimal digits", s)
	}

	if id == uuid.Nil() {
		return uuid.UUID{}, fmt.Errorf(`%q is not a canonical RFC 9562 UUID: the "nil" UUID is not supported`, s)
	}

	return id, nil
}
