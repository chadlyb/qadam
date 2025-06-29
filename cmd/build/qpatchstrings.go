package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/chadlyb/qadam/shared"
)

// Expects lines in the format:
// BEGIN-END:"STRING" ; COMMENT
//
// BEGIN-END: is followed by a hex offset, and a colon.
// STRING is a quoted string.
// ; COMMENT is optional (and ignored)
const lineRegexSrc = `^\s*(?:0x)?(?P<begin>[0-9a-fA-F]+)\s*-\s*(?:0x)?(?P<end>[0-9a-fA-F]+)\s*:\s*"(?P<string>(?:[^"\\]|\\"|\\n|\\\\|\\t|\\r)*)"\s*(?:;.*)?$`

var lineRegex = regexp.MustCompile(lineRegexSrc)

func handleLine(data []byte, line string) error {
	matches := lineRegex.FindStringSubmatch(line)
	if len(matches) != 4 {
		return fmt.Errorf("line didn't match expected format: %v", line)
	}

	patchBegin, err := strconv.ParseUint(matches[1], 16, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse begin offset: %w", err)
	}
	patchEnd, err := strconv.ParseUint(matches[2], 16, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse end offset: %w", err)
	}
	patchBytes, err := shared.FromString(matches[3])
	if err != nil {
		return fmt.Errorf("couldn't translate string: %w", err)
	}

	patchLen := uint64(len(patchBytes))
	if patchLen+1 > patchEnd-patchBegin {
		return fmt.Errorf("ignoring too-long string (%v > %v)", patchLen+1, patchEnd-patchBegin)
	}

	for i := uint64(0); i != patchLen; i++ {
		data[patchBegin+i] = patchBytes[i]
	}
	data[patchBegin+patchLen] = 0
	return nil
}

// qpatchStringsFromReader processes data from io.Reader and patch data from io.Reader, writing results to io.Writer
func qpatchStringsFromReader(srcReader io.Reader, destWriter io.Writer, patchReader io.Reader) error {
	// Read source data
	data, err := io.ReadAll(srcReader)
	if err != nil {
		return fmt.Errorf("couldn't read source data: %w", err)
	}

	// Parse patch data, line by line
	scanner := bufio.NewScanner(patchReader)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		err := handleLine(data, line)
		if err != nil {
			fmt.Printf("warning: ignored line %v due to error: %v\n", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("couldn't scan patch data: %w", err)
	}

	// Write to output
	_, err = destWriter.Write(data)
	if err != nil {
		return fmt.Errorf("couldn't write output: %w", err)
	}

	return nil
}

// qpatchStrings is the convenience function that maintains the original file path interface
func qpatchStrings(srcPath string, destPath string, patchPath string) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't read source file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	// Open patch file
	patchFile, err := os.Open(patchPath)
	if err != nil {
		return fmt.Errorf("couldn't open patch file %s: %w", patchPath, err)
	}
	defer patchFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	return qpatchStringsFromReader(srcFile, destFile, patchFile)
}
