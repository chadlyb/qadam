package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

// fixgame patches the expected size of two files in the game executable to the size of the files in the source directory
func fixgame(srcPath string, destPath string, filesToGetLengthFrom ...string) error {
	// Default patch offsets for TEXTS.FIL and RESOURCE.FIL
	defaultPatchOffsets := []int{0x0001A706, 0x0001A6E6}

	if len(filesToGetLengthFrom) != len(defaultPatchOffsets) {
		return fmt.Errorf("wrong number of source paths, expecting %v got %v", len(defaultPatchOffsets), len(filesToGetLengthFrom))
	}

	patchValues := make([]uint32, len(defaultPatchOffsets))
	for i := 0; i != len(defaultPatchOffsets); i++ {
		fileInfo, err := os.Lstat(filesToGetLengthFrom[i])
		if err != nil {
			return fmt.Errorf("fixgame could not get size of file '%v': %w", filesToGetLengthFrom[i], err)
		}
		patchValues[i] = uint32(fileInfo.Size())
	}

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", destPath, err)
	}
	defer destFile.Close()

	return fixgameFromReader(srcFile, destFile, patchValues, defaultPatchOffsets)
}
