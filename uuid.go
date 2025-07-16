package dogma

import (
	"fmt"
)

func normalizeUUID(id string) (string, error) {
	if len(id) != 36 {
		return "", fmt.Errorf("%q is not a canonical RFC 9562 UUID: expected 36 characters", id)
	}

	var normalized [36]byte
	isNil := true

	for i := 0; i < 36; i++ {
		c := id[i]
		normalized[i] = c

		switch i {
		case 8, 13, 18, 23:
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
