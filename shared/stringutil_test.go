package shared

import (
	"bytes"
	"testing"
)

func TestUnescapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "simple string",
			input:    "hello world",
			expected: "hello world",
			hasError: false,
		},
		{
			name:     "newline escape",
			input:    "hello\\nworld",
			expected: "hello\nworld",
			hasError: false,
		},
		{
			name:     "tab escape",
			input:    "hello\\tworld",
			expected: "hello\tworld",
			hasError: false,
		},
		{
			name:     "carriage return escape",
			input:    "hello\\rworld",
			expected: "hello\rworld",
			hasError: false,
		},
		{
			name:     "backslash escape",
			input:    "hello\\\\world",
			expected: "hello\\world",
			hasError: false,
		},
		{
			name:     "quote escape",
			input:    "hello\\\"world",
			expected: "hello\"world",
			hasError: false,
		},
		{
			name:     "hex escape",
			input:    "hello\\x20world",
			expected: "hello world",
			hasError: false,
		},
		{
			name:     "multiple escapes",
			input:    "\\n\\t\\r\\\"\\\\",
			expected: "\n\t\r\"\\",
			hasError: false,
		},
		{
			name:     "trailing backslash",
			input:    "hello\\",
			expected: "",
			hasError: true,
		},
		{
			name:     "invalid hex escape",
			input:    "hello\\xzworld",
			expected: "",
			hasError: true,
		},
		{
			name:     "incomplete hex escape",
			input:    "hello\\x1",
			expected: "",
			hasError: true,
		},
		{
			name:     "unknown escape",
			input:    "hello\\zworld",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnescapeString(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple bytes",
			input:    []byte{72, 101, 108, 108, 111}, // "Hello"
			expected: "Hello",
		},
		{
			name:     "with quotes",
			input:    []byte{34, 72, 101, 108, 108, 111, 34}, // "\"Hello\""
			expected: "\\\"Hello\\\"",
		},
		{
			name:     "with newlines",
			input:    []byte{72, 101, 108, 108, 111, 10, 119, 111, 114, 108, 100}, // "Hello\nworld"
			expected: "Hello\\nworld",
		},
		{
			name:     "with tabs",
			input:    []byte{72, 101, 108, 108, 111, 9, 119, 111, 114, 108, 100}, // "Hello\tworld"
			expected: "Hello\\tworld",
		},
		{
			name:     "with backslashes",
			input:    []byte{72, 101, 108, 108, 111, 92, 119, 111, 114, 108, 100}, // "Hello\\world"
			expected: "Hello\\\\world",
		},
		{
			name:     "czech characters",
			input:    []byte{0x8E, 0x8F, 0x90}, // Some Czech characters
			expected: string([]rune{CharsetRunes[0x8E], CharsetRunes[0x8F], CharsetRunes[0x90]}),
		},
		{
			name:     "control characters",
			input:    []byte{0, 1, 2, 3}, // Control characters
			expected: string([]rune{CharsetRunes[0], CharsetRunes[1], CharsetRunes[2], CharsetRunes[3]}),
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		hasError bool
	}{
		{
			name:     "simple string",
			input:    "Hello",
			expected: []byte{72, 101, 108, 108, 111},
			hasError: false,
		},
		{
			name:     "with escaped quotes",
			input:    "\\\"Hello\\\"",
			expected: []byte{34, 72, 101, 108, 108, 111, 34},
			hasError: false,
		},
		{
			name:     "with escaped newlines",
			input:    "Hello\\nworld",
			expected: []byte{72, 101, 108, 108, 111, 10, 119, 111, 114, 108, 100},
			hasError: false,
		},
		{
			name:     "with escaped tabs",
			input:    "Hello\\tworld",
			expected: []byte{72, 101, 108, 108, 111, 9, 119, 111, 114, 108, 100},
			hasError: false,
		},
		{
			name:     "with escaped backslashes",
			input:    "Hello\\\\world",
			expected: []byte{72, 101, 108, 108, 111, 92, 119, 111, 114, 108, 100},
			hasError: false,
		},
		{
			name:     "czech characters",
			input:    string([]rune{CharsetRunes[0x8E], CharsetRunes[0x8F], CharsetRunes[0x90]}),
			expected: []byte{0x8E, 0x8F, 0x90},
			hasError: false,
		},
		{
			name:     "control characters",
			input:    string([]rune{CharsetRunes[0], CharsetRunes[1], CharsetRunes[2]}),
			expected: []byte{0, 1, 2},
			hasError: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
			hasError: false,
		},
		{
			name:     "unrecognized character",
			input:    "Hello\u0000world", // NUL character not in charset
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid escape sequence",
			input:    "Hello\\zworld",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromString(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !bytes.Equal(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestCharsetMapping(t *testing.T) {
	// Test that charset mapping is consistent
	t.Run("charset consistency", func(t *testing.T) {
		// Convert CharsetString to runes for proper comparison
		charsetStringRunes := []rune(CharsetString)

		if len(CharsetRunes) != len(charsetStringRunes) {
			t.Errorf("CharsetRunes length (%d) != CharsetString runes length (%d)", len(CharsetRunes), len(charsetStringRunes))
		}

		for i, r := range CharsetRunes {
			if r != charsetStringRunes[i] {
				t.Errorf("Mismatch at position %d: CharsetRunes[%d] = %c, CharsetString rune[%d] = %c", i, i, r, i, charsetStringRunes[i])
			}
		}
	})

	t.Run("charset map consistency", func(t *testing.T) {
		for i, r := range CharsetRunes {
			if byteVal, exists := CharsetMapToByte[r]; !exists {
				t.Errorf("Character %c at position %d not found in CharsetMapToByte", r, i)
			} else if byteVal != byte(i) {
				t.Errorf("Character %c at position %d maps to %d, expected %d", r, i, byteVal, i)
			}
		}
	})

	t.Run("round trip conversion", func(t *testing.T) {
		testBytes := []byte{0, 1, 2, 10, 32, 65, 97, 128, 255}
		for _, b := range testBytes {
			if b < byte(len(CharsetRunes)) {
				r := CharsetRunes[b]
				if mappedByte, exists := CharsetMapToByte[r]; !exists {
					t.Errorf("Byte %d -> rune %c not found in map", b, r)
				} else if mappedByte != b {
					t.Errorf("Byte %d -> rune %c -> byte %d (expected %d)", b, r, mappedByte, b)
				}
			}
		}
	})
}

func TestSpecialUnicodeMappings(t *testing.T) {
	tests := []struct {
		name     string
		unicode  rune
		expected byte
	}{
		{"en dash", '\u2013', CharsetMapToByte['-']},
		{"em dash", '\u2014', CharsetMapToByte['-']},
		{"left single quote", '\u2018', CharsetMapToByte['\'']},
		{"right single quote", '\u2019', CharsetMapToByte['\'']},
		{"left double quote", '\u201A', CharsetMapToByte['"']},
		{"right double quote", '\u201B', CharsetMapToByte['"']},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if mappedByte, exists := CharsetMapToByte[tt.unicode]; !exists {
				t.Errorf("Unicode character %c not found in map", tt.unicode)
			} else if mappedByte != tt.expected {
				t.Errorf("Unicode character %c maps to %d, expected %d", tt.unicode, mappedByte, tt.expected)
			}
		})
	}
}

func TestControlCharacterMappings(t *testing.T) {
	tests := []struct {
		name     string
		control  rune
		expected byte
	}{
		{"carriage return", '\r', '\r'},
		{"newline", '\n', '\n'},
		{"tab", '\t', '\t'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if mappedByte, exists := CharsetMapToByte[tt.control]; !exists {
				t.Errorf("Control character %c not found in map", tt.control)
			} else if mappedByte != tt.expected {
				t.Errorf("Control character %c maps to %d, expected %d", tt.control, mappedByte, tt.expected)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"simple ascii", []byte("Hello, World!")},
		{"with spaces", []byte("Hello World")},
		{"with numbers", []byte("12345")},
		{"mixed case", []byte("HelloWorld")},
		{"czech text", []byte{0x8E, 0x8F, 0x90, 0x91, 0x92}}, // Some Czech characters
		{"control chars", []byte{0, 1, 2, 10, 13, 9}},        // Control characters
		{"special chars", []byte{34, 39, 92, 45}},            // Quotes, backslash, dash
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert bytes to string
			str := ToString(tt.input)

			// Convert string back to bytes
			result, err := FromString(str)
			if err != nil {
				t.Errorf("Round trip conversion failed: %v", err)
				return
			}

			// Compare original and result
			if !bytes.Equal(tt.input, result) {
				t.Errorf("Round trip conversion failed:\nOriginal: %v\nResult:   %v", tt.input, result)
			}
		})
	}
}

func TestFindNextValidString(t *testing.T) {
	tests := []struct {
		name          string
		data          []byte
		startPos      int
		endPos        int
		catchAll      bool
		expectedStart int
		expectedEnd   int
		expectedFound bool
		description   string
	}{
		{
			name:          "empty data",
			data:          []byte{},
			startPos:      0,
			endPos:        0,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   0,
			expectedFound: false,
			description:   "Empty data should return not found with endPos",
		},
		{
			name:          "too short data",
			data:          []byte{0x41, 0x42},
			startPos:      0,
			endPos:        2,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   2,
			expectedFound: false,
			description:   "Data shorter than MinStringLength should return not found with endPos",
		},
		{
			name:          "valid string with null terminator",
			data:          []byte{0x41, 0x42, 0x43, 0x00, 0x44, 0x45},
			startPos:      0,
			endPos:        6,
			catchAll:      false,
			expectedStart: 0,
			expectedEnd:   3,
			expectedFound: true,
			description:   "Simple valid string should return correct start/end offsets",
		},
		{
			name:          "valid string at end",
			data:          []byte{0x00, 0x00, 0x41, 0x42, 0x43},
			startPos:      0,
			endPos:        5,
			catchAll:      false,
			expectedStart: 2,
			expectedEnd:   5,
			expectedFound: true,
			description:   "String at end should return correct offsets",
		},
		{
			name:          "catch all mode - any non-empty string",
			data:          []byte{0x00, 0x01, 0x02, 0x00, 0x41, 0x42},
			startPos:      0,
			endPos:        6,
			catchAll:      true,
			expectedStart: 1,
			expectedEnd:   3,
			expectedFound: true,
			description:   "Catch-all mode should accept any non-empty string",
		},
		{
			name:          "no valid string start",
			data:          []byte{0x00, 0x01, 0x02, 0x00, 0x00, 0x00},
			startPos:      0,
			endPos:        6,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   6,
			expectedFound: false,
			description:   "No valid string start should return not found with endPos",
		},
		{
			name:          "multiple strings - find first",
			data:          []byte{0x00, 'H', 'e', 'l', 'l', 'o', 0x00, 'W', 'o', 'r', 'l', 'd', 0x00},
			startPos:      0,
			endPos:        13,
			catchAll:      false,
			expectedStart: 1,
			expectedEnd:   6,
			expectedFound: true,
			description:   "Should find first valid string in multiple strings",
		},
		{
			name:          "multiple strings - start from middle",
			data:          []byte{0x00, 'H', 'e', 'l', 'l', 'o', 0x00, 'W', 'o', 'r', 'l', 'd', 0x00},
			startPos:      7,
			endPos:        13,
			catchAll:      false,
			expectedStart: 7,
			expectedEnd:   12,
			expectedFound: true,
			description:   "Should find string when starting from middle position",
		},
		{
			name:          "string with garbage prefix",
			data:          []byte{0x00, 0x01, 0x02, 'H', 'e', 'l', 'l', 'o', 0x00},
			startPos:      0,
			endPos:        9,
			catchAll:      false,
			expectedStart: 3,
			expectedEnd:   8,
			expectedFound: true,
			description:   "Should find string after garbage prefix",
		},
		{
			name:          "string with garbage suffix",
			data:          []byte{'H', 'e', 'l', 'l', 'o', 0x00, 0x01, 0x02, 0x03},
			startPos:      0,
			endPos:        9,
			catchAll:      false,
			expectedStart: 0,
			expectedEnd:   5,
			expectedFound: true,
			description:   "Should find string before garbage suffix",
		},
		{
			name:          "short valid string in catch-all mode",
			data:          []byte{0x00, 'H', 'i', 0x00},
			startPos:      0,
			endPos:        4,
			catchAll:      true,
			expectedStart: 1,
			expectedEnd:   3,
			expectedFound: true,
			description:   "Catch-all should accept short strings",
		},
		{
			name:          "short invalid string in conservative mode",
			data:          []byte{0x00, 'H', 'i', 0x00},
			startPos:      0,
			endPos:        4,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   4,
			expectedFound: false,
			description:   "Conservative mode should reject short strings",
		},
		{
			name:          "null bytes only",
			data:          []byte{0x00, 0x00, 0x00, 0x00},
			startPos:      0,
			endPos:        4,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   4,
			expectedFound: false,
			description:   "Null bytes only should return not found",
		},
		{
			name:          "null bytes only in catch-all",
			data:          []byte{0x00, 0x00, 0x00, 0x00},
			startPos:      0,
			endPos:        4,
			catchAll:      true,
			expectedStart: -1,
			expectedEnd:   4,
			expectedFound: false,
			description:   "Null bytes only should return not found even in catch-all",
		},
		{
			name:          "partial scan - start in middle",
			data:          []byte{'H', 'e', 'l', 'l', 'o', 0x00, 'W', 'o', 'r', 'l', 'd', 0x00},
			startPos:      3,
			endPos:        12,
			catchAll:      false,
			expectedStart: 6,
			expectedEnd:   11,
			expectedFound: true,
			description:   "Partial scan should find next valid string",
		},
		{
			name:          "realistic game data scenario",
			data:          []byte{0x00, 0x00, 'B', 'a', 'r', 'l', 'a', 'n', 'd', ' ', 'C', '+', '+', ' ', '-', ' ', 'C', 'o', 'p', 'y', 'r', 'i', 'g', 'h', 't', ' ', '9', '9', '9', '1', ' ', 'B', 'a', 'r', 'l', 'a', 'n', 'd', ' ', 'I', 'n', 't', 'l', '.', 0x00, 0x00, 0x00},
			startPos:      0,
			endPos:        47,
			catchAll:      false,
			expectedStart: 2,
			expectedEnd:   44,
			expectedFound: true,
			description:   "Realistic Barland copyright string scenario",
		},
		{
			name:          "edge case - start at end",
			data:          []byte{'H', 'e', 'l', 'l', 'o', 0x00},
			startPos:      6,
			endPos:        6,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   6,
			expectedFound: false,
			description:   "Starting at end should return not found",
		},
		{
			name:          "edge case - start past end",
			data:          []byte{'H', 'e', 'l', 'l', 'o', 0x00},
			startPos:      10,
			endPos:        6,
			catchAll:      false,
			expectedStart: -1,
			expectedEnd:   6,
			expectedFound: false,
			description:   "Starting past end should return not found",
		},
		{
			name:          "consecutive null terminators",
			data:          []byte{'H', 'e', 'l', 'l', 'o', 0x00, 0x00, 'W', 'o', 'r', 'l', 'd', 0x00},
			startPos:      0,
			endPos:        13,
			catchAll:      false,
			expectedStart: 0,
			expectedEnd:   5,
			expectedFound: true,
			description:   "Should handle consecutive null terminators correctly",
		},
		{
			name:          "string without null terminator",
			data:          []byte{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'},
			startPos:      0,
			endPos:        11,
			catchAll:      false,
			expectedStart: 0,
			expectedEnd:   11,
			expectedFound: true,
			description:   "Should handle string without null terminator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultStart, resultEnd, resultFound := FindNextValidString(tt.data, tt.startPos, tt.endPos, tt.catchAll)

			if resultStart != tt.expectedStart {
				t.Errorf("FindNextValidString() start = %v, want %v (%s)", resultStart, tt.expectedStart, tt.description)
			}
			if resultEnd != tt.expectedEnd {
				t.Errorf("FindNextValidString() end = %v, want %v (%s)", resultEnd, tt.expectedEnd, tt.description)
			}
			if resultFound != tt.expectedFound {
				t.Errorf("FindNextValidString() found = %v, want %v (%s)", resultFound, tt.expectedFound, tt.description)
			}

			// Additional verification: if found, verify the string content
			if resultFound && resultStart >= 0 && resultEnd > resultStart {
				actualString := string(tt.data[resultStart:resultEnd])
				if tt.catchAll {
					// In catch-all mode, any non-empty string is acceptable
					if len(actualString) == 0 {
						t.Errorf("Found empty string in catch-all mode (%s)", tt.description)
					}
				} else {
					// In conservative mode, verify it starts with acceptable character
					if len(actualString) > 0 && !IsAcceptableStringStart(tt.data[resultStart]) {
						t.Errorf("String doesn't start with acceptable character: %q (%s)", actualString, tt.description)
					}
				}
			}
		})
	}
}
