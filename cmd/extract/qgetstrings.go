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

	fmt.Printf("DEBUG: Read %d bytes from %s\n", len(data), srcPath)

	outStream, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't create file '%v': %w", destPath, err)
	}
	defer outStream.Close()

	// Crawl for strings, output OFFSET LENGTH "<string>"

	foundBorland := false
	totalStrings := 0
	potentialStrings := 0
	acceptedStrings := 0

	stringBegin := 0
	at := 0
	end := len(data)
	for at != end {
		if data[at] == 0 {
			if at-stringBegin > 2 { // Only process strings longer than 2 chars
				potentialStrings++
				s := shared.ToString(data[stringBegin:at])

				// Debug: Check for Borland string
				if !foundBorland && strings.Contains(s, "Borland") {
					foundBorland = true
					fmt.Printf("DEBUG: Found Borland string at offset %08x: \"%s\"\n", stringBegin, s)
				}

				// Debug: Show some potential strings
				if potentialStrings <= 10 {
					fmt.Printf("DEBUG: Potential string %d at %08x: \"%s\" (len=%d)\n",
						potentialStrings, stringBegin, s, at-stringBegin)
				}

				if foundBorland {
					totalStrings++
					isLikely := shared.IsLikelyHumanLanguage(data[stringBegin:at])
					if isLikely {
						acceptedStrings++
						fmt.Fprintf(outStream, "%08x-%08x: \"%v\"\n", stringBegin, at+1, s)

						// Debug: Show accepted strings
						if acceptedStrings <= 5 {
							fmt.Printf("DEBUG: ACCEPTED string %d: \"%s\"\n", acceptedStrings, s)
						}
					} else {
						// Debug: Show rejected strings
						if totalStrings <= 10 {
							fmt.Printf("DEBUG: REJECTED string %d: \"%s\"\n", totalStrings, s)
						}
					}
				}
			}
			stringBegin = at + 1
		}

		at++
	}

	// Handle final string if no null terminator
	if stringBegin != at && at-stringBegin > 2 {
		potentialStrings++
		if foundBorland {
			totalStrings++
			if shared.IsLikelyHumanLanguage(data[stringBegin:at]) {
				acceptedStrings++
				fmt.Fprintf(outStream, "%08x-%08x: \"%v\" NO_NUL\n", stringBegin, at, shared.ToString(data[stringBegin:at]))
			}
		}
	}

	fmt.Printf("DEBUG: Summary for %s:\n", srcPath)
	fmt.Printf("  - Total potential strings: %d\n", potentialStrings)
	fmt.Printf("  - Strings after Borland: %d\n", totalStrings)
	fmt.Printf("  - Accepted as Czech: %d\n", acceptedStrings)
	fmt.Printf("  - Found Borland marker: %v\n", foundBorland)

	return nil
}
