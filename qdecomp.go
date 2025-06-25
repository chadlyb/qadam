package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func qdecomp(inputFile string, outputFile string) error {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Failed to read file '%v': %w\n", inputFile, err)
	}

	numEntries := int(data[0])
	offsets := make([]int, numEntries+1)
	for i := 0; i != numEntries+1; i++ {
		offsets[i] = int(data[1+i*3+2])<<16 + int(data[1+i*3+1])<<8 + int(data[1+i*3+0])
	}

	outStream, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Failed to create file '%v': %w\n", outputFile, err)
	}
	defer outStream.Close()

	// Todo: Validate offsets[numEntries] == len(data)

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
					fmt.Fprintf(outStream, "%c", charsetRunes[b])
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
}
