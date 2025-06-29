package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestQDecompFromReader(t *testing.T) {
	// Test data representing a simple .FIL file with 1 section
	// Format: [numEntries][offset1][offset2][data...]
	// Header: 1 byte + 2 * 3 bytes = 7 bytes
	testData := []byte{
		0x01,             // 1 entry
		0x07, 0x00, 0x00, // offset 1: 7 (start of data section)
		0x0C, 0x00, 0x00, // offset 2: 12 (end of data)
		// Data section (5 bytes)
		0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello" (hex data)
	}

	// Create input reader and output writer
	reader := bytes.NewReader(testData)
	var writer bytes.Buffer

	// Run the function
	err := qdecompFromReader(reader, &writer)
	if err != nil {
		t.Fatalf("qdecompFromReader failed: %v", err)
	}

	output := writer.String()
	t.Logf("Output: %s", output)

	// Verify results
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have at least one line
	if len(lines) == 0 {
		t.Error("No output was generated")
	}

	// Should contain SECTION 0
	foundSection := false
	for _, line := range lines {
		if strings.Contains(line, "SECTION 0") {
			foundSection = true
			break
		}
	}

	if !foundSection {
		t.Error("Expected to find 'SECTION 0' in output")
	}

	// Should contain hex data
	foundHex := false
	for _, line := range lines {
		if strings.Contains(line, "[48 65 6C 6C 6F]") {
			foundHex = true
			break
		}
	}

	if !foundHex {
		t.Error("Expected to find hex data '[48 65 6C 6C 6F]' in output")
	}
}

func TestQDecompFromReaderWithString(t *testing.T) {
	// Test data with a string section
	// Format: [numEntries][offset1][offset2][data...]
	// Header: 1 byte + 2 * 3 bytes = 7 bytes
	// Section: [5 hex bytes][zero-terminated obfuscated string]
	// "Hi" in charset: H=0x48, i=0x69; obfuscated with +0x31: 0x48+0x31=0x79, 0x69+0x31=0x9A
	testData := []byte{
		0x01,             // 1 entry
		0x07, 0x00, 0x00, // offset 1: 7 (start of data section)
		0x0F, 0x00, 0x00, // offset 2: 15 (end of data)
		// Data section: 5 hex bytes, then obfuscated string "Hi" + NUL
		0x01, 0x02, 0x03, 0x04, 0x05, // 5 hex bytes
		0x79, 0x9A, 0x00, // "Hi" obfuscated, NUL-terminated
	}

	// Create input reader and output writer
	reader := bytes.NewReader(testData)
	var writer bytes.Buffer

	// Run the function
	err := qdecompFromReader(reader, &writer)
	if err != nil {
		t.Fatalf("qdecompFromReader failed: %v", err)
	}

	output := writer.String()
	t.Logf("Output: %s", output)

	// Verify results
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have at least one line
	if len(lines) == 0 {
		t.Error("No output was generated")
	}

	// Should contain SECTION 0
	foundSection := false
	for _, line := range lines {
		if strings.Contains(line, "SECTION 0") {
			foundSection = true
			break
		}
	}
	if !foundSection {
		t.Error("Expected to find 'SECTION 0' in output")
	}

	// Should contain the hex run
	foundHex := false
	for _, line := range lines {
		if strings.Contains(line, "[01 02 03 04 05]") {
			foundHex = true
			break
		}
	}
	if !foundHex {
		t.Error("Expected to find '[01 02 03 04 05]' in output")
	}

	// Should contain the string "Hi"
	foundString := false
	for _, line := range lines {
		if strings.Contains(line, "\"Hi\"") {
			foundString = true
			break
		}
	}
	if !foundString {
		t.Error("Expected to find '\"Hi\"' in output")
	}
}
