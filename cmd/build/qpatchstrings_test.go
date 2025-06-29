package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestQPatchStringsFromReader(t *testing.T) {
	// Test source data (simple executable-like data)
	srcData := []byte{
		// Some header bytes
		0x4D, 0x5A, 0x90, 0x00, 0x03, 0x00, 0x00, 0x00,
		// Some data that will be patched
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x00, // "Hello World\0"
		// More data
		0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x65, 0x73, 0x74, 0x00, // "This is test\0"
	}

	// Test patch data - fix the ranges to match the string lengths
	patchData := "00000008-0000000D:\"Ahoj\"\n0000000E-00000016:\"Svět\"\n"

	// Create readers and writer
	srcReader := bytes.NewReader(srcData)
	patchReader := strings.NewReader(patchData)
	var destWriter bytes.Buffer

	// Run the function
	err := qpatchStringsFromReader(srcReader, &destWriter, patchReader)
	if err != nil {
		t.Fatalf("qpatchStringsFromReader failed: %v", err)
	}

	output := destWriter.Bytes()
	t.Logf("Output length: %d bytes", len(output))

	// Verify the patches were applied correctly
	// Original: "Hello World\0" at offset 8-20
	// Patched: "Ahoj\0" at offset 8-13
	expectedPatch1 := []byte{0x41, 0x68, 0x6F, 0x6A, 0x00} // "Ahoj\0" in charset
	if !bytes.Equal(output[8:13], expectedPatch1) {
		t.Errorf("First patch failed. Expected %v, got %v", expectedPatch1, output[8:13])
	}

	// Original: "This is test\0" at offset 21-34
	// Patched: "Svět\0" at offset 14-19
	expectedPatch2 := []byte{0x53, 0x76, 0xD8, 0x74, 0x00} // "Svět\0" in charset
	if !bytes.Equal(output[14:19], expectedPatch2) {
		t.Errorf("Second patch failed. Expected %v, got %v", expectedPatch2, output[14:19])
	}

	// Verify other data remains unchanged
	if output[0] != 0x4D || output[1] != 0x5A {
		t.Error("Header data was modified unexpectedly")
	}
}

func TestQPatchStringsFromReaderWithInvalidPatch(t *testing.T) {
	// Test source data
	srcData := []byte{
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00, // "Hello\0"
	}

	// Test patch data with invalid format (should be ignored with warning)
	patchData := "invalid format line\n00000000-00000005:\"Ahoj\"\n"

	// Create readers and writer
	srcReader := bytes.NewReader(srcData)
	patchReader := strings.NewReader(patchData)
	var destWriter bytes.Buffer

	// Run the function - should not fail, just warn about invalid line
	err := qpatchStringsFromReader(srcReader, &destWriter, patchReader)
	if err != nil {
		t.Fatalf("qpatchStringsFromReader failed: %v", err)
	}

	output := destWriter.Bytes()
	t.Logf("Output length: %d bytes", len(output))

	// Verify the valid patch was applied
	expectedPatch := []byte{0x41, 0x68, 0x6F, 0x6A, 0x00} // "Ahoj\0" in charset
	if !bytes.Equal(output[0:5], expectedPatch) {
		t.Errorf("Valid patch failed. Expected %v, got %v", expectedPatch, output[0:5])
	}
}
