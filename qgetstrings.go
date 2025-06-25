package main

import (
	"fmt"
	"os"
)

func qgetStrings(srcPath string, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("Error reading %s: %w", srcPath, err)
	}

	outStream, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Failed to create file '%v': %w\n", destPath, err)
	}
	defer outStream.Close()
	// Crawl for strings, output OFFSET LENGTH "<string>"

	stringBegin := 0
	at := 0
	end := len(data)
	for at != end {
		if data[at] == 0 {
			if heuristicIsHumanString(data[stringBegin:at]) {
				fmt.Fprintf(outStream, "%08x-%08x: \"%v\"\n", stringBegin, at+1, toString(data[stringBegin:at]))
			}
			stringBegin = at + 1
		}

		at++
	}
	if stringBegin != at {
		if heuristicIsHumanString(data[stringBegin:at]) {
			fmt.Fprintf(outStream, "%08x-%08x: \"%v\" NO_NUL\n", stringBegin, at, toString(data[stringBegin:at]))
		}
	}
	return nil
}
