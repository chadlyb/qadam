package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chadlyb/qadam/shared"
)

func readInt24(data []byte, offset int) int32 {
	return int32(data[offset]) | int32(data[offset+1])<<8 | int32(data[offset+2])<<16
}

// qdecompFromReader processes data from an io.Reader and writes results to an io.Writer
func qdecompFromReader(reader io.Reader, writer io.Writer) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	if len(data) < 1 {
		return fmt.Errorf("data is empty")
	}

	numEntries := int(data[0])
	if numEntries < 0 || numEntries > 1000 {
		return fmt.Errorf("invalid number of entries: %v", numEntries)
	}

	offsets := make([]int, numEntries+1)
	for i := 0; i != numEntries+1; i++ {
		offsets[i] = int(readInt24(data, 1+i*3))
		if offsets[i] < 0 || offsets[i] > len(data) {
			return fmt.Errorf("offset %v is out of bounds", offsets[i])
		}
	}

	if offsets[numEntries] != len(data) {
		return fmt.Errorf("last offset %v does not match data size %v", offsets[numEntries], len(data))
	}

	hexRun := 0
	stringRun := 0

	const HEX_RUN_LIMIT = 5
	begin := 1 + (numEntries+1)*3
	sectionEnd := begin
	for at := begin; ; at++ {
		if at == sectionEnd {
			if stringRun > 0 {
				// This is bad...
				fmt.Fprintf(writer, "\" NO_NUL\n") // I guess we just hack this and can handle it in recompiler...
			} else if hexRun > 0 {
				fmt.Fprintf(writer, "]\n")
			}
			for i := 0; i != numEntries; i++ {
				if offsets[i] == at {
					sectionEnd = offsets[i+1]
					fmt.Fprintf(writer, "SECTION %v\n", i)
				}
			}
			if at == len(data) {
				break
			}
			hexRun = 0
			stringRun = 0
		}

		if stringRun > 0 || hexRun == HEX_RUN_LIMIT {
			if stringRun == 0 && hexRun == HEX_RUN_LIMIT {
				fmt.Fprintf(writer, "] \"")
				hexRun = 0
			}
			if data[at] == 0x00 {
				fmt.Fprintf(writer, "\"\n")
				stringRun = 0
			} else {
				b := data[at] - 0x31 // unobfuscate.
				switch {
				case b == '\n':
					fmt.Fprintf(writer, "\\n")
				case b == '\t':
					fmt.Fprintf(writer, "\\t")
				case b == '\\':
					fmt.Fprintf(writer, "\\\\")
				case b == '"':
					fmt.Fprintf(writer, "\\\"")
				default:
					fmt.Fprintf(writer, "%c", shared.CharsetRunes[b])
				}
				stringRun++
			}
		} else {
			if hexRun == 0 {
				fmt.Fprintf(writer, "[")
			} else {
				fmt.Fprintf(writer, " ")
			}
			fmt.Fprintf(writer, "%02X", data[at])
			hexRun++
		}
	}

	fmt.Fprintf(writer, "\n")
	return nil
}

// qdecomp is the convenience function that maintains the original file path interface
func qdecomp(inputFile string, outputFile string) error {
	// Open input file
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file '%v': %w", inputFile, err)
	}
	defer input.Close()

	// Create output file
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file '%v': %w", outputFile, err)
	}
	defer output.Close()

	return qdecompFromReader(input, output)
}
