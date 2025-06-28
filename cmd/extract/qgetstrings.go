package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chadlyb/qadam/shared"
)

func qgetStrings(srcPath string, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't read %s: %w", srcPath, err)
	}

	outStream, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't create file '%v': %w", destPath, err)
	}
	defer outStream.Close()

	// Crawl for strings, output OFFSET LENGTH "<string>"

	foundBorland := false

	stringBegin := 0
	at := 0
	end := len(data)
	for at != end {
		if data[at] == 0 {
			s := shared.ToString(data[stringBegin:at])
			if !foundBorland && strings.Contains(s, "Borland") {
				foundBorland = true
			}
			if foundBorland && at-stringBegin > 2 {
				if shared.IsLikelyCzechString(data[stringBegin:at]) {
					fmt.Fprintf(outStream, "%08x-%08x: \"%v\"\n", stringBegin, at+1, shared.ToString(data[stringBegin:at]))
				}
			}
			stringBegin = at + 1
		}

		at++
	}
	if stringBegin != at {
		if shared.IsLikelyCzechString(data[stringBegin:at]) {
			fmt.Fprintf(outStream, "%08x-%08x: \"%v\" NO_NUL\n", stringBegin, at, shared.ToString(data[stringBegin:at]))
		}
	}
	return nil
}
