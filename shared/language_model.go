package shared

// Valid human-readable characters (Czech includes all English characters)
var validHumanChars = map[byte]bool{
	// Basic Latin letters
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true, 'y': true, 'z': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,
	// Space and common punctuation
	' ': true, '.': true, ',': true, '!': true, '?': true, ':': true, ';': true, '-': true, '_': true, '(': true, ')': true, '"': true, '\'': true,
	// Czech accented characters (lowercase)
	0xA0: true, // á
	0xA1: true, // í
	0xA3: true, // ú
	0xA7: true, // ž
	0x9F: true, // č
	0x85: true, // ů
	0x82: true, // é
	0xE7: true, // š
	0xD8: true, // ě
	0xE5: true, // ň
	0xEC: true, // ý
	0xFD: true, // ř
	0x9C: true, // ť
	0xD4: true, // ď
	// Czech accented characters (uppercase)
	0xB5: true, // Á
	0xD6: true, // Í
	0xE9: true, // Ú
	0xA6: true, // Ž
	0xAC: true, // Č
	0xDE: true, // Ů
	0x90: true, // É
	0xE6: true, // Š
	0xB7: true, // Ě
	0xD5: true, // Ň
	0xED: true, // Ý
	0xFC: true, // Ř
	0x9B: true, // Ť
	0xD2: true, // Ď
}

// IsLikelyHumanLanguage checks if at least 2/3 of the string consists of valid human-readable characters
func IsLikelyHumanLanguage(bytes []byte) bool {
	const minLength = 3
	if len(bytes) < minLength {
		return false
	}

	validCount := 0
	for _, b := range bytes {
		if validHumanChars[b] {
			validCount++
		}
	}

	// Require at least 2/3 of characters to be valid
	threshold := float64(len(bytes)) * 2.0 / 3.0
	return float64(validCount) >= threshold
}
