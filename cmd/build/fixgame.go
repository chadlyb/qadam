package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func fixgame(srcPath string, destPath string, srcPaths ...string) error {
	patch_locs := []int{0x0001A706, 0x0001A6E6}

	if len(srcPaths) != len(patch_locs) {
		return fmt.Errorf("wrong number of source paths, expecting %v got %v", len(patch_locs), len(srcPaths))
	}

	patch_values := make([]uint32, len(patch_locs))
	for i := 0; i != len(patch_locs); i++ {
		fileInfo, err := os.Lstat(srcPaths[i])
		if err != nil {
			return fmt.Errorf("fixgame could not get size of file '%v': %w", srcPaths[i], err)
		}
		patch_values[i] = uint32(fileInfo.Size())
	}

	// Read input file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", srcPath, err)
	}

	for i := 0; i != len(patch_locs); i++ {
		// Check bounds
		if patch_locs[i]+2 > len(data) {
			return fmt.Errorf("error: offset 0x%X out of bounds for GAME.EXE (size %d bytes)", patch_locs[i], len(data))
		}

		// Patch the value in-place
		binary.LittleEndian.PutUint32(data[patch_locs[i]:], uint32(patch_values[i]))
	}

	// Write to output file
	err = os.WriteFile(destPath, data, 0755)
	if err != nil {
		return fmt.Errorf("error writing %s: %w", destPath, err)
	}

	//fmt.Printf("Patched %s -> %s", srcPath, destPath)
	return nil
}
