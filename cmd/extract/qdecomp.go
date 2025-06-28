package main

import (
	"fmt"
	"os"

	"github.com/chadlyb/qadam/shared"
)

func readInt24(data []byte, offset int) int32 {
	return int32(data[offset]) | int32(data[offset+1])<<8 | int32(data[offset+2])<<16
}

func qdecomp(inputFile string, outputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file '%v': %w", inputFile, err)
	}

	if len(data) < 1 {
		return fmt.Errorf("file %v is empty", inputFile)
	}

	numEntries := int(data[0])
	if numEntries < 0 || numEntries > 1000 {
		return fmt.Errorf("invalid number of entries: %v", numEntries)
	}

	offsets := make([]int, numEntries+1)
	for i := 0; i != numEntries+1; i++ {
		offsets[i] = int(readInt24(data, 1+i*3))
		if offsets[i] < 0 || offsets[i] > len(data) {
			return fmt.Errorf("offset %v is out of bounds for file %v", offsets[i], inputFile)
		}
	}

	outStream, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file '%v': %w", outputFile, err)
	}
	defer outStream.Close()

	if offsets[numEntries] != len(data) {
		return fmt.Errorf("last offset %v does not match file size %v", offsets[numEntries], len(data))
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
				fmt.Fprintf(outStream, "\" NO_NUL\n") // I guess we just hack this and can handle it in recompiler...
			} else if hexRun > 0 {
				fmt.Fprintf(outStream, "]\n")
			}
			for i := 0; i != numEntries; i++ {
				if offsets[i] == at {
					sectionEnd = offsets[i+1]
					fmt.Fprintf(outStream, "SECTION %v\n", i)
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
				fmt.Fprintf(outStream, "] \"")
				hexRun = 0
			}
			if data[at] == 0x00 {
				fmt.Fprintf(outStream, "\"\n")
				stringRun = 0
			} else {
				b := data[at] - 0x31 // unobfuscate.
				switch {
				case b == '\n':
					fmt.Fprintf(outStream, "\\n")
				case b == '\t':
					fmt.Fprintf(outStream, "\\t")
				case b == '\\':
					fmt.Fprintf(outStream, "\\\\")
				case b == '"':
					fmt.Fprintf(outStream, "\\\"")
				default:
					fmt.Fprintf(outStream, "%c", shared.CharsetRunes[b])
				}
				stringRun++
			}
		} else {
			if hexRun == 0 {
				fmt.Fprintf(outStream, "[")
			} else {
				fmt.Fprintf(outStream, " ")
			}
			fmt.Fprintf(outStream, "%02X", data[at])
			hexRun++
		}
	}

	fmt.Fprintf(outStream, "\n")
	return nil
}
