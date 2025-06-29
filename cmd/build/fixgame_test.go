package main

import (
	"bytes"
	"encoding/binary"
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

func TestFixgameInPlacePatching(t *testing.T) {
	// Test that we can patch a file "in place" using the buffer approach
	// This simulates the scenario where we read a file, patch it, and write it back

	// Create test data that simulates a game executable
	testData := make([]byte, 0x1A710)
	for i := range testData {
		testData[i] = byte(i & 0xFF)
	}

	// Set some initial values at the patch offsets
	binary.LittleEndian.PutUint32(testData[0x1A706:], 0x12345678) // TEXTS.FIL size
	binary.LittleEndian.PutUint32(testData[0x1A6E6:], 0x87654321) // RESOURCE.FIL size

	// Create a buffer with the test data
	srcBuffer := bytes.NewReader(testData)
	var dstBuffer bytes.Buffer

	// Patch the data
	fileSizes := []uint32{0x1000, 0x2000} // New sizes
	patchOffsets := []int{0x0001A706, 0x0001A6E6}

	err := fixgameFromReader(srcBuffer, &dstBuffer, fileSizes, patchOffsets)
	if err != nil {
		t.Fatalf("fixgameFromReader failed: %v", err)
	}

	result := dstBuffer.Bytes()

	// Verify the patches were applied correctly
	patchedTextsSize := binary.LittleEndian.Uint32(result[0x1A706:])
	patchedResourceSize := binary.LittleEndian.Uint32(result[0x1A6E6:])

	if patchedTextsSize != 0x1000 {
		t.Errorf("TEXTS.FIL size patch failed. Expected 0x1000, got 0x%X", patchedTextsSize)
	}

	if patchedResourceSize != 0x2000 {
		t.Errorf("RESOURCE.FIL size patch failed. Expected 0x2000, got 0x%X", patchedResourceSize)
	}

	// Verify other data remains unchanged
	if result[0] != 0x00 || result[1] != 0x01 {
		t.Error("Non-patch data was modified unexpectedly")
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
