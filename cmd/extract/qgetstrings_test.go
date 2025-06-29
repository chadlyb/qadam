package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chadlyb/qadam/shared"
)

func TestQGetStrings(t *testing.T) {
	// Create a temporary test file with some sample data
	testData := []byte{
		// Some random bytes
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00, // "Hello\0"
		// Borland marker
		0x42, 0x6F, 0x72, 0x6C, 0x61, 0x6E, 0x64, 0x00, // "Borland\0"
		// Some Czech text
		0x41, 0x68, 0x6F, 0x6A, 0x20, 0x73, 0x76, 0xD8, 0x74, 0x65, 0x00, // "Ahoj světe\0" (using correct charset)
		// Some English text
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x00, // "Hello world\0"
		// More Czech text
		0x54, 0x72, 0x65, 0x73, 0x74, 0x6E, 0x69, 0x63, 0x65, 0x20, 0x6E, 0x61, 0xE7, 0x65, 0x00, // "Trestnice naše\0"
	}

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.exe")
	outputFile := filepath.Join(tempDir, "output.txt")

	// Write test data
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run string extraction
	err = qgetStrings(testFile, outputFile)
	if err != nil {
		t.Fatalf("qgetStrings failed: %v", err)
	}

	// Read output and verify
	output, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(output)
	t.Logf("Output: %s", outputStr)

	// Should find at least some Czech strings
	if len(outputStr) == 0 {
		t.Error("No strings were extracted")
	}

	// Should contain the Czech strings we added
	if !containsString(outputStr, "Ahoj světe") {
		t.Error("Expected to find 'Ahoj světe' in output")
	}

	if !containsString(outputStr, "Trestnice naše") {
		t.Error("Expected to find 'Trestnice naše' in output")
	}
}

// TestQGetStringsFromReader demonstrates the improved testability with buffers
func TestQGetStringsFromReader(t *testing.T) {
	// Test data with Borland marker and Czech strings
	testData := []byte{
		// Borland marker
		0x42, 0x6F, 0x72, 0x6C, 0x61, 0x6E, 0x64, 0x00, // "Borland\0"
		// Czech text
		0x41, 0x68, 0x6F, 0x6A, 0x20, 0x73, 0x76, 0xD8, 0x74, 0x65, 0x00, // "Ahoj světe\0"
		// English text
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x00, // "Hello world\0"
		// Numbers (should be rejected)
		0x31, 0x32, 0x33, 0x34, 0x35, 0x00, // "12345\0"
	}

	// Create input reader and output writer
	reader := bytes.NewReader(testData)
	var writer bytes.Buffer

	// Run the function
	err := qgetStringsFromReader(reader, &writer)
	if err != nil {
		t.Fatalf("qgetStringsFromReader failed: %v", err)
	}

	output := writer.String()
	t.Logf("Output: %s", output)

	// Verify results
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have extracted strings
	if len(lines) == 0 {
		t.Error("No strings were extracted")
	}

	// Should contain Czech and English strings
	foundCzech := false
	foundEnglish := false
	foundNumbers := false

	for _, line := range lines {
		if strings.Contains(line, "Ahoj světe") {
			foundCzech = true
		}
		if strings.Contains(line, "Hello world") {
			foundEnglish = true
		}
		if strings.Contains(line, "12345") {
			foundNumbers = true
		}
	}

	if !foundCzech {
		t.Error("Expected to find 'Ahoj světe' in output")
	}
	if !foundEnglish {
		t.Error("Expected to find 'Hello world' in output")
	}
	if foundNumbers {
		t.Error("Expected '12345' to be rejected")
	}
}

func TestIsLikelyHumanStringWithRealData(t *testing.T) {
	// Test with some real Czech strings
	testCases := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Game menu", "Hra", true},
		{"Settings", "Nastavení", true},
		{"Exit", "Konec", true},
		{"Save", "Uložit", true},
		{"Load", "Načíst", true},
		{"English menu", "Game", true},         // English is valid human language
		{"English settings", "Settings", true}, // English is valid human language
		{"Numbers", "123", false},
		{"Mixed", "Game123", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := shared.FromString(tc.text)
			if err != nil {
				t.Skipf("Skipping test due to conversion error: %v", err)
				return
			}

			result := shared.IsLikelyHumanLanguage(bytes)
			if result != tc.expected {
				t.Errorf("IsLikelyHumanLanguage(\"%s\") = %v, expected %v", tc.text, result, tc.expected)
			}
		})
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		(haystack == needle ||
			len(haystack) > len(needle) &&
				(haystack[:len(needle)] == needle ||
					haystack[len(haystack)-len(needle):] == needle ||
					containsSubstring(haystack, needle)))
}

func containsSubstring(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
