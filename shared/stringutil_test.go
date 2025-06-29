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
		name           string
		data           []byte
		startPos       int
		endPos         int
		expectedString string
		expectedNewPos int
		expectedFound  bool
	}{
		{
			name:           "simple_valid_string",
			data:           []byte{0x00, 0x00, 'H', 'e', 'l', 'l', 'o', 0x00, 0x00},
			startPos:       2,
			endPos:         9,
			expectedString: "Hello",
			expectedNewPos: 8,
			expectedFound:  true,
		},
		{
			name:           "czech_string",
			data:           []byte{0x00, 0xA0, 0xE7, 0x82, 0x00, 0x00}, // ášé
			startPos:       1,
			endPos:         6,
			expectedString: "ášé",
			expectedNewPos: 5,
			expectedFound:  true,
		},
		{
			name:           "string_starts_with_bad_character",
			data:           []byte{0x00, ' ', 'H', 'e', 'l', 'l', 'o', 0x00},
			startPos:       1,
			endPos:         8,
			expectedString: "Hello",
			expectedNewPos: 8,
			expectedFound:  true,
		},
		{
			name:           "string_too_short",
			data:           []byte{0x00, 'H', 'i', 0x00},
			startPos:       1,
			endPos:         4,
			expectedString: "",
			expectedNewPos: 4,
			expectedFound:  false,
		},
		{
			name:           "no_null_terminator",
			data:           []byte{0x00, 'H', 'e', 'l', 'l', 'o'},
			startPos:       1,
			endPos:         6,
			expectedString: "Hello",
			expectedNewPos: 7,
			expectedFound:  true,
		},
		{
			name:           "range_too_small",
			data:           []byte{'H', 'i'},
			startPos:       0,
			endPos:         2,
			expectedString: "",
			expectedNewPos: 2,
			expectedFound:  false,
		},
		{
			name:           "start_at_end",
			data:           []byte{'H', 'e', 'l', 'l', 'o'},
			startPos:       5,
			endPos:         5,
			expectedString: "",
			expectedNewPos: 5,
			expectedFound:  false,
		},
		{
			name:           "multiple_strings_find_first",
			data:           []byte{'H', 'i', 0x00, 'H', 'e', 'l', 'l', 'o', 0x00},
			startPos:       0,
			endPos:         9,
			expectedString: "Hello",
			expectedNewPos: 9,
			expectedFound:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultString, resultNewPos, resultFound := FindNextValidString(tt.data, tt.startPos, tt.endPos)

			if resultString != tt.expectedString {
				t.Errorf("FindNextValidString() string = %v, want %v", resultString, tt.expectedString)
			}
			if resultNewPos != tt.expectedNewPos {
				t.Errorf("FindNextValidString() newPos = %v, want %v", resultNewPos, tt.expectedNewPos)
			}
			if resultFound != tt.expectedFound {
				t.Errorf("FindNextValidString() found = %v, want %v", resultFound, tt.expectedFound)
			}
		})
	}
}
