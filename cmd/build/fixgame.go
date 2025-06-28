package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// fixgame patches the expected size of two files in the game executable to the size of the files in the source directory
func fixgame(srcPath string, destPath string, filesToGetLengthFrom ...string) error {
	patchOffsets := []int{0x0001A706, 0x0001A6E6} // These correspond to the offsets of the expected sizes of TEXTS.FIL and RESOURCE.FIL in the game executable

	if len(filesToGetLengthFrom) != len(patchOffsets) {
		return fmt.Errorf("wrong number of source paths, expecting %v got %v", len(patchOffsets), len(filesToGetLengthFrom))
	}

	patchValues := make([]uint32, len(patchOffsets))
	for i := 0; i != len(patchOffsets); i++ {
		fileInfo, err := os.Lstat(filesToGetLengthFrom[i])
		if err != nil {
			return fmt.Errorf("fixgame could not get size of file '%v': %w", filesToGetLengthFrom[i], err)
		}
		patchValues[i] = uint32(fileInfo.Size())
	}

	// Read input file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", srcPath, err)
	}

	for i := 0; i != len(patchOffsets); i++ {
		// Check bounds
		if patchOffsets[i]+4 > len(data) {
			return fmt.Errorf("error: offset 0x%X out of bounds for %v (size %d bytes)", patchOffsets[i], srcPath, len(data))
		}

		// Patch the value in-place
		binary.LittleEndian.PutUint32(data[patchOffsets[i]:], uint32(patchValues[i]))
	}

	// Write to output file
	err = os.WriteFile(destPath, data, 0755)
	if err != nil {
		return fmt.Errorf("error writing %s: %w", destPath, err)
	}

	//fmt.Printf("Patched %s -> %s", srcPath, destPath)
	return nil
}
