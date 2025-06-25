package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// BEGIN-END:"STRING" ; COMMENT
const lineRegexSrc = `^\s*(?:0x)?(?P<begin>[0-9a-fA-F]+)\s*-\s*(?:0x)?(?P<end>[0-9a-fA-F]+)\s*:\s*"(?P<string>(?:[^"\\]|\\"|\\n|\\\\|\\t|\\r)*)"\s*(?:;.*)?$`

var lineRegex = regexp.MustCompile(lineRegexSrc)

func handleLine(data []byte, line string) error {
	matches := lineRegex.FindStringSubmatch(line)
	//fmt.Println(matches)
	if len(matches) != 4 {
		return fmt.Errorf("Line didn't match expected format")
	}

	patchBegin, err := strconv.ParseUint(matches[1], 16, 64)
	if err != nil {
		return fmt.Errorf("Couldn't parse begin offset (%w)", err)
	}
	patchEnd, err := strconv.ParseUint(matches[2], 16, 64)
	if err != nil {
		return fmt.Errorf("Couldn't parse end offset (%w)", err)
	}
	patchBytes, err := fromString(matches[3])
	if err != nil {
		return fmt.Errorf("Couldn't translate string (%w)", err)
	}

	len := uint64(len(patchBytes))
	if len+1 > patchEnd-patchBegin {
		return fmt.Errorf("Ignoring too-long string (%v > %v)", len+1, patchEnd-patchBegin)
	}

	for i := uint64(0); i != len; i++ {
		data[patchBegin+i] = patchBytes[i]
	}
	data[patchBegin+len] = 0
	return nil
}

func qpatchStrings(srcPath string, destPath string, patchPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("Error reading source file %s: %w", srcPath, err)
	}

	// Parse patch path, line by line.
	patchFile, err := os.Open(patchPath)
	if err != nil {
		return fmt.Errorf("Error opening patch file %s: %w", patchPath, err)
	}
	defer patchFile.Close()

	scanner := bufio.NewScanner(patchFile)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		err := handleLine(data, line)
		if err != nil {
			fmt.Printf("Ignored line %v due to error: %v\n", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error while scanning patch file %s: %w", patchPath, err)
	}

	// Write to output file
	err = os.WriteFile(destPath, data, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", destPath, err)
		os.Exit(1)
	}

	fmt.Printf("Patched %s -> %s", srcPath, destPath)
	return nil
}

func main() {
	err := realMain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
