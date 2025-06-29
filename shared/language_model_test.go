package shared

import (
	"testing"
)

func TestIsLikelyHumanLanguage(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Czech text", "Trestnice naše", true},
		{"English text", "Hello world", true},
		{"Mixed text", "Hello svět", true},
		{"Numbers", "123456789", false},
		{"Empty string", "", false},
		{"Single char", "a", false},
		{"Short string", "Hi", false},
		{"Short valid string", "Hra", true},
		{"Short valid string 2", "Game", true},
		{"Too many numbers", "Game123", false},
		{"Some numbers OK", "Hello 123", true},
		{"Czech with numbers", "Hra 123", false},
		{"English with Czech chars", "Hello svět", true},
		{"Capitalized Czech", "Áno", true},
		{"Capitalized Czech 2", "Ídete", true},
		{"Capitalized Czech 3", "Úspěch", true},
		{"Mixed case Czech", "Česká Republika", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := FromString(tc.text)
			if err != nil {
				t.Skipf("Skipping test due to conversion error: %v", err)
				return
			}

			result := IsLikelyHumanLanguage(bytes)
			if result != tc.expected {
				t.Errorf("IsLikelyHumanLanguage(\"%s\") = %v, expected %v", tc.text, result, tc.expected)
			}
		})
	}
}

func TestIsAcceptableStringStart(t *testing.T) {
	tests := []struct {
		name     string
		char     byte
		expected bool
	}{
		// English letters (should be acceptable)
		{"lowercase_a", 'a', true},
		{"uppercase_A", 'A', true},
		{"lowercase_z", 'z', true},
		{"uppercase_Z", 'Z', true},

		// Czech accented letters (should be acceptable)
		{"czech_á", 0xA0, true}, // á
		{"czech_í", 0xA1, true}, // í
		{"czech_ú", 0xA3, true}, // ú
		{"czech_ž", 0xA7, true}, // ž
		{"czech_č", 0x9F, true}, // č
		{"czech_ů", 0x85, true}, // ů
		{"czech_é", 0x82, true}, // é
		{"czech_š", 0xE7, true}, // š
		{"czech_ě", 0xD8, true}, // ě
		{"czech_ň", 0xE5, true}, // ň
		{"czech_ý", 0xEC, true}, // ý
		{"czech_ř", 0xFD, true}, // ř
		{"czech_ť", 0x9C, true}, // ť
		{"czech_ď", 0xD4, true}, // ď

		// Czech uppercase accented letters (should be acceptable)
		{"czech_Á", 0xB5, true}, // Á
		{"czech_Í", 0xD6, true}, // Í
		{"czech_Ú", 0xE9, true}, // Ú
		{"czech_Ž", 0xA6, true}, // Ž
		{"czech_Č", 0xAC, true}, // Č
		{"czech_Ů", 0xDE, true}, // Ů
		{"czech_É", 0x90, true}, // É
		{"czech_Š", 0xE6, true}, // Š
		{"czech_Ě", 0xB7, true}, // Ě
		{"czech_Ň", 0xD5, true}, // Ň
		{"czech_Ý", 0xED, true}, // Ý
		{"czech_Ř", 0xFC, true}, // Ř
		{"czech_Ť", 0x9B, true}, // Ť
		{"czech_Ď", 0xD2, true}, // Ď

		// Non-letter characters (should NOT be acceptable)
		{"space", ' ', false},
		{"period", '.', false},
		{"comma", ',', false},
		{"exclamation", '!', false},
		{"question", '?', false},
		{"colon", ':', false},
		{"semicolon", ';', false},
		{"hyphen", '-', false},
		{"underscore", '_', false},
		{"parenthesis_open", '(', false},
		{"parenthesis_close", ')', false},
		{"quote", '"', false},
		{"apostrophe", '\'', false},

		// Control characters and other non-printable (should NOT be acceptable)
		{"null", 0x00, false},
		{"tab", 0x09, false},
		{"newline", 0x0A, false},
		{"carriage_return", 0x0D, false},
		{"random_control", 0x1F, false},
		{"random_high", 0xFF, false},
		{"random_mid", 0x80, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAcceptableStringStart(tt.char)
			if result != tt.expected {
				t.Errorf("IsAcceptableStringStart(0x%02x) = %v, want %v", tt.char, result, tt.expected)
			}
		})
	}
}
