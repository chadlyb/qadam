package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

func main() {
	patch_locs := []int{0x0001A706, 0x0001A6E6}

	if len(os.Args) != len(patch_locs)+3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source EXE> <dest EXE> <texts.fil size (decimal)> <resource.fil size (decimal)>\n", os.Args[0])
		os.Exit(1)
	}

	srcPath := os.Args[1]
	destPath := os.Args[2]

	patch_values := make([]uint32, len(patch_locs))
	for i := 0; i != len(patch_locs); i++ {
		// Parse decimal input
		val, err := strconv.ParseUint(os.Args[i+3], 10, 32)
		if err != nil {
                       fileInfo, err := os.Lstat(os.Args[i+3])
                       if err != nil {
                               fmt.Fprintf(os.Stderr, "Invalid decimal uint32 value and also not could not open as file that exists '%v': %v\n", os.Args[i+3], err)
                               os.Exit(1)
                       }
                       val = uint64(fileInfo.Size())
		}
		patch_values[i] = uint32(val)
	}

	// Read input file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", srcPath, err)
		os.Exit(1)
	}

	for i := 0; i != len(patch_locs); i++ {
		// Check bounds
		if patch_locs[i]+2 > len(data) {
			fmt.Fprintf(os.Stderr, "Error: offset 0x%X out of bounds for GAME.EXE (size %d bytes)\n", patch_locs[i], len(data))
			os.Exit(1)
		}

		// Patch the value in-place
		binary.LittleEndian.PutUint32(data[patch_locs[i]:], uint32(patch_values[i]))
	}

	// Write to output file
	err = os.WriteFile(destPath, data, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", destPath, err)
		os.Exit(1)
	}

	fmt.Printf("Patched %s -> %s", srcPath, destPath)
}
