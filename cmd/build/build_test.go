package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestGameExecutablePatching(t *testing.T) {
	// Test that we can patch a game executable with correct file sizes
	// This simulates the scenario where we read a file, patch it, and write it back

	// Create test data that simulates a game executable
	testData := make([]byte, 0x1A710)
	for i := range testData {
		testData[i] = byte(i & 0xFF)
	}

	// Set some initial values at the patch offsets
	binary.LittleEndian.PutUint32(testData[0x1A706:], 0x12345678) // TEXTS.FIL size
	binary.LittleEndian.PutUint32(testData[0x1A6E6:], 0x87654321) // RESOURCE.FIL size

	// Create a temporary file
	tempDir := t.TempDir()
	gameExePath := filepath.Join(tempDir, "GAME.EXE")

	// Write test data to file
	err := os.WriteFile(gameExePath, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Simulate the patching process
	gameExeData, err := os.ReadFile(gameExePath)
	if err != nil {
		t.Fatalf("Failed to read game executable: %v", err)
	}

	// Patch the file sizes directly in the buffer
	const textsFilOffset = 0x0001A706
	const resourceFilOffset = 0x0001A6E6

	// Check bounds
	if textsFilOffset+4 > len(gameExeData) {
		t.Fatalf("TEXTS.FIL offset 0x%X out of bounds (size %d bytes)", textsFilOffset, len(gameExeData))
	}
	if resourceFilOffset+4 > len(gameExeData) {
		t.Fatalf("RESOURCE.FIL offset 0x%X out of bounds (size %d bytes)", resourceFilOffset, len(gameExeData))
	}

	// Patch the values in-place
	binary.LittleEndian.PutUint32(gameExeData[textsFilOffset:], 0x1000)    // New TEXTS.FIL size
	binary.LittleEndian.PutUint32(gameExeData[resourceFilOffset:], 0x2000) // New RESOURCE.FIL size

	// Write the patched data back to the file
	err = os.WriteFile(gameExePath, gameExeData, 0644)
	if err != nil {
		t.Fatalf("Failed to write patched game executable: %v", err)
	}

	// Read back and verify the patches were applied correctly
	patchedData, err := os.ReadFile(gameExePath)
	if err != nil {
		t.Fatalf("Failed to read patched file: %v", err)
	}

	patchedTextsSize := binary.LittleEndian.Uint32(patchedData[0x1A706:])
	patchedResourceSize := binary.LittleEndian.Uint32(patchedData[0x1A6E6:])

	if patchedTextsSize != 0x1000 {
		t.Errorf("TEXTS.FIL size patch failed. Expected 0x1000, got 0x%X", patchedTextsSize)
	}

	if patchedResourceSize != 0x2000 {
		t.Errorf("RESOURCE.FIL size patch failed. Expected 0x2000, got 0x%X", patchedResourceSize)
	}

	// Verify other data remains unchanged
	if patchedData[0] != 0x00 || patchedData[1] != 0x01 {
		t.Error("Non-patch data was modified unexpectedly")
	}
}
