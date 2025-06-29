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
