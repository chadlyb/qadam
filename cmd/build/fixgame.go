package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// fixgameFromReader patches the expected size of files in the game executable data
func fixgameFromReader(srcReader io.Reader, destWriter io.Writer, fileSizes []uint32, patchOffsets []int) error {
	if len(fileSizes) != len(patchOffsets) {
		return fmt.Errorf("wrong number of file sizes, expecting %v got %v", len(patchOffsets), len(fileSizes))
	}

	// Read input data
	data, err := io.ReadAll(srcReader)
	if err != nil {
		return fmt.Errorf("error reading source data: %w", err)
	}

	for i := 0; i != len(patchOffsets); i++ {
		// Check bounds
		if patchOffsets[i]+4 > len(data) {
			return fmt.Errorf("error: offset 0x%X out of bounds (size %d bytes)", patchOffsets[i], len(data))
		}

		// Patch the value in-place
		binary.LittleEndian.PutUint32(data[patchOffsets[i]:], fileSizes[i])
	}

	// Write to output
	_, err = destWriter.Write(data)
	if err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}
