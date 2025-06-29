package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chadlyb/qadam/shared"
)

// qgetStringsFromReader processes data from an io.Reader and writes results to an io.Writer
func qgetStringsFromReader(reader io.Reader, writer io.Writer, catchAll bool) error {
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
	skippedBadStart := 0

	pos := 0
	end := len(data)

	for pos < end {
		// Find the next valid string starting from current position
		stringStart, stringEnd, found := shared.FindNextValidString(data, pos, end, catchAll)

		if !found {
			// No valid string found, move to next position and continue
			pos = stringEnd + 1
			continue
		}

		// Calculate the original start position for this string
		pos = stringEnd + 1

		potentialStrings++

		// Generate the string content from the range
		stringContent := shared.ToString(data[stringStart:stringEnd])

		// Debug: Check for Borland string
		if !foundBorland && strings.Contains(stringContent, "Borland") {
			foundBorland = true
			if debugMode {
				fmt.Printf("DEBUG: Found Borland string at offset %08x: \"%s\"\n", stringStart, stringContent)
			}
		}

		// Debug: Show all potential strings when in debug mode
		if debugMode {
			fmt.Printf("DEBUG: Potential string %d at %08x: \"%s\" (len=%d)\n",
				potentialStrings, stringStart, stringContent, len(stringContent))
		}

		// In catch-all mode, process all strings regardless of Borland or length
		if catchAll || foundBorland {
			totalStrings++
			if catchAll {
				acceptedStrings++
				fmt.Fprintf(writer, "%08x-%08x: \"%v\"\n", stringStart, stringEnd+1, stringContent)
				if debugMode {
					fmt.Printf("DEBUG: ACCEPTED string %d: \"%s\"\n", acceptedStrings, stringContent)
				}
				continue
			}
			// Conservative mode: Convert string back to bytes for language detection
			stringBytes, err := shared.FromString(stringContent)
			if err != nil {
				if debugMode {
					fmt.Printf("DEBUG: Error converting string back to bytes: %v\n", err)
				}
				continue
			}
			isLikely := shared.IsLikelyHumanLanguage(stringBytes)
			if isLikely {
				acceptedStrings++
				fmt.Fprintf(writer, "%08x-%08x: \"%v\"\n", stringStart, stringEnd+1, stringContent)
				if debugMode {
					fmt.Printf("DEBUG: ACCEPTED string %d: \"%s\"\n", acceptedStrings, stringContent)
				}
			} else {
				if debugMode {
					fmt.Printf("DEBUG: REJECTED string %d: \"%s\"\n", totalStrings, stringContent)
				}
			}
		}
	}

	if debugMode {
		fmt.Printf("DEBUG: Summary:\n")
		fmt.Printf("  - Total potential strings: %d\n", potentialStrings)
		fmt.Printf("  - Strings after Borland: %d\n", totalStrings)
		fmt.Printf("  - Accepted as Czech: %d\n", acceptedStrings)
		fmt.Printf("  - Skipped bad start: %d\n", skippedBadStart)
		fmt.Printf("  - Found Borland marker: %v\n", foundBorland)
		fmt.Printf("  - Catch-all mode: %v\n", catchAll)
	}

	return nil
}

// qgetStrings extracts strings from a file and writes them to another file
func qgetStrings(srcPath, destPath string, allStrings bool) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't create destination file: %w", err)
	}
	defer destFile.Close()

	return qgetStringsFromReader(srcFile, destFile, allStrings)
}
