package dogma

import (
	"fmt"
)

func normalizeUUID(id string) (string, error) {
	const size = 36

	if len(id) != size {
		return "", fmt.Errorf("%q is not a canonical RFC 9562 UUID: expected 36 characters", id)
	}

	var normalized [size]byte
	isNil := true

	for i := range size {
		c := id[i]
		normalized[i] = c

		switch i {
		case 8, 13, 18, 23: // indexes of hyphens
			if c != '-' {
				return "", fmt.Errorf("%q is not a canonical RFC 9562 UUID: expected hyphen at position %d", id, i)
			}
		default:
			switch {
			case c == '0':
				// ok
			case c >= '1' && c <= '9':
				isNil = false
			case c >= 'a' && c <= 'f':
				isNil = false
			case c >= 'A' && c <= 'F':
				isNil = false
				normalized[i] += 'a' - 'A' // convert to lowercase
			default:
				return "", fmt.Errorf("%q is not a canonical RFC 9562 UUID: expected hex digit at position %d", id, i)
			}
		}
	}

	if isNil {
		return "", fmt.Errorf(`%q is not a canonical RFC 9562 UUID: the "nil" UUID is not supported`, id)
	}

	return string(normalized[:]), nil
}
