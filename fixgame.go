package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <uint16 decimal value>\n", os.Args[0])
		os.Exit(1)
	}

	// Parse decimal input
	val, err := strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid decimal uint16 value: %v\n", err)
		os.Exit(1)
	}

	const offset = 0x0001A706

	// Read input file
	data, err := os.ReadFile("GAME.EXE")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading GAME.EXE: %v\n", err)
		os.Exit(1)
	}

	// Check bounds
	if offset+2 > len(data) {
		fmt.Fprintf(os.Stderr, "Error: Offset 0x%X out of bounds for GAME.EXE (size %d bytes)\n", offset, len(data))
		os.Exit(1)
	}

	// Patch the value in-place
	binary.LittleEndian.PutUint16(data[offset:], uint16(val))

	// Write to output file
	err = os.WriteFile("GAME2.EXE", data, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing GAME2.EXE: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Patched GAME.EXE -> GAME2.EXE with value %d at offset 0x%X\n", val, offset)
}

