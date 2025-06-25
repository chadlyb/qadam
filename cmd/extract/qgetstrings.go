package main

import (
	"fmt"
	"os"

	"github.com/chadlyb/qadam/shared"
)

func qgetStrings(srcPath string, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", srcPath, err)
	}

	outStream, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file '%v': %w", destPath, err)
	}
	defer outStream.Close()
	// Crawl for strings, output OFFSET LENGTH "<string>"

	stringBegin := 0
	at := 0
	end := len(data)
	for at != end {
		if data[at] == 0 {
			if shared.HeuristicIsHumanString(data[stringBegin:at]) {
				fmt.Fprintf(outStream, "%08x-%08x: \"%v\"\n", stringBegin, at+1, shared.ToString(data[stringBegin:at]))
			}
			stringBegin = at + 1
		}

		at++
	}
	if stringBegin != at {
		if shared.HeuristicIsHumanString(data[stringBegin:at]) {
			fmt.Fprintf(outStream, "%08x-%08x: \"%v\" NO_NUL\n", stringBegin, at, shared.ToString(data[stringBegin:at]))
		}
	}
	return nil
}
