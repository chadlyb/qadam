package main

import (
	"bytes"
	"testing"
)

func TestFixgameFromReader(t *testing.T) {
	// Test data: 16 bytes with two 4-byte values to patch
	input := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Patch offsets: at positions 4 and 8
	patchOffsets := []int{4, 8}
	fileSizes := []uint32{0x12345678, 0x87654321}

	var output bytes.Buffer
	err := fixgameFromReader(bytes.NewReader(input), &output, fileSizes, patchOffsets)

	if err != nil {
		t.Fatalf("fixgameFromReader failed: %v", err)
	}

	result := output.Bytes()

	// Check that the values were patched correctly
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x78, 0x56, 0x34, 0x12, 0x21, 0x43, 0x65, 0x87, 0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(result, expected) {
		t.Errorf("fixgameFromReader output mismatch\ngot:  %x\nwant: %x", result, expected)
	}
}

func TestFixgameFromReaderWithWrongFileCount(t *testing.T) {
	// Test source data
	srcData := make([]byte, 0x1A710)

	// Test with wrong number of file sizes (should fail)
	fileSizes := []uint32{0x1000} // Only 1 size, but we need 2
	patchOffsets := []int{0x0001A706, 0x0001A6E6}

	// Create reader and writer
	srcReader := bytes.NewReader(srcData)
	var destWriter bytes.Buffer

	// Run the function - should fail
	err := fixgameFromReader(srcReader, &destWriter, fileSizes, patchOffsets)
	if err == nil {
		t.Error("Expected error for wrong number of file sizes, but got none")
	}

	expectedError := "wrong number of file sizes, expecting 2 got 1"
	if !contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func TestFixgameFromReaderWithOutOfBounds(t *testing.T) {
	// Test source data that's too small
	srcData := make([]byte, 0x1A700) // Too small for our patch offsets

	// Test file sizes
	fileSizes := []uint32{0x1000, 0x2000}
	patchOffsets := []int{0x0001A706, 0x0001A6E6}

	// Create reader and writer
	srcReader := bytes.NewReader(srcData)
	var destWriter bytes.Buffer

	// Run the function - should fail
	err := fixgameFromReader(srcReader, &destWriter, fileSizes, patchOffsets)
	if err == nil {
		t.Error("Expected error for out of bounds offset, but got none")
	}

	expectedError := "offset 0x1A706 out of bounds"
	if !contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
