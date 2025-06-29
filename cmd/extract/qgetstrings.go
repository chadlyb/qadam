package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chadlyb/qadam/shared"
)

// qgetStringsFromReader processes data from an io.Reader and writes results to an io.Writer
func qgetStringsFromReader(reader io.Reader, writer io.Writer) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("couldn't read data: %w", err)
	}

	if debugMode {
		fmt.Printf("DEBUG: Read %d bytes\n", len(data))
	}

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
					if debugMode {
						fmt.Printf("DEBUG: Found Borland string at offset %08x: \"%s\"\n", stringBegin, s)
					}
				}

				// Debug: Show all potential strings when in debug mode
				if debugMode {
					fmt.Printf("DEBUG: Potential string %d at %08x: \"%s\" (len=%d)\n",
						potentialStrings, stringBegin, s, at-stringBegin)
				}

				if foundBorland {
					totalStrings++
					isLikely := shared.IsLikelyHumanLanguage(data[stringBegin:at])
					if isLikely {
						acceptedStrings++
						fmt.Fprintf(writer, "%08x-%08x: \"%v\"\n", stringBegin, at+1, s)

						// Debug: Show all accepted strings when in debug mode
						if debugMode {
							fmt.Printf("DEBUG: ACCEPTED string %d: \"%s\"\n", acceptedStrings, s)
						}
					} else {
						// Debug: Show all rejected strings when in debug mode
						if debugMode {
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
				fmt.Fprintf(writer, "%08x-%08x: \"%v\" NO_NUL\n", stringBegin, at, shared.ToString(data[stringBegin:at]))
			}
		}
	}

	if debugMode {
		fmt.Printf("DEBUG: Summary:\n")
		fmt.Printf("  - Total potential strings: %d\n", potentialStrings)
		fmt.Printf("  - Strings after Borland: %d\n", totalStrings)
		fmt.Printf("  - Accepted as Czech: %d\n", acceptedStrings)
		fmt.Printf("  - Found Borland marker: %v\n", foundBorland)
	}

	return nil
}

// qgetStrings is the convenience function that maintains the original file path interface
func qgetStrings(srcPath string, destPath string) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't open %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't create file '%v': %w", destPath, err)
	}
	defer destFile.Close()

	if debugMode {
		fmt.Printf("DEBUG: Processing %s -> %s\n", srcPath, destPath)
	}

	return qgetStringsFromReader(srcFile, destFile)
}
