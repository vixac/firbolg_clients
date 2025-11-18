package bullet_stl

import "fmt"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
const base = int64(36)

// BulletIdIntToaasci converts a non-negative int64 to a fixed-length base36 string.
// Returns an error if intValue is out of range for the chosen idSize.
func BulletIdIntToaasci(intValue int64, idSize int) (string, error) {
	if intValue < 0 {
		return "", fmt.Errorf("value must be non-negative: %d", intValue)
	}

	// Maximum encodable value is 36^idSize - 1
	maxVal := int64(1)
	for i := 0; i < idSize; i++ {
		maxVal *= base
	}

	if intValue >= maxVal {
		return "", fmt.Errorf("value %d is too large for idSize %d (max %d)", intValue, idSize, maxVal-1)
	}

	// Convert to base36
	buf := make([]byte, 0, idSize)
	v := intValue

	if v == 0 {
		buf = append(buf, '0')
	} else {
		for v > 0 {
			d := v % base
			v /= base
			buf = append(buf, alphabet[d])
		}
	}

	// Reverse buffer
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	// Left-pad with '0' to fixed width
	if len(buf) < idSize {
		pad := make([]byte, idSize-len(buf))
		for i := range pad {
			pad[i] = '0'
		}
		buf = append(pad, buf...)
	}

	return string(buf), nil
}

// AasciBulletIdToInt decodes a base36 ID string back into an int64.
func AasciBulletIdToInt(aasci string) (int64, error) {
	var value int64 = 0

	for _, c := range aasci {
		var digit int64 = -1

		switch {
		case c >= '0' && c <= '9':
			digit = int64(c - '0')
		case c >= 'a' && c <= 'z':
			digit = int64(c-'a') + 10
		default:
			return 0, fmt.Errorf("invalid character '%c' in id string", c)
		}

		value = value*base + digit
	}

	return value, nil
}
