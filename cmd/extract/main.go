package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chadlyb/qadam/shared"
)

// Version will be set by the linker during build
var version = "dev"

// Global debug flag
var debugMode = false

func extract(srcPath string, outputDir string, allStrings bool) error {
	// Use provided output directory or default to ../extracted relative to source
	if outputDir == "" {
		outputDir = filepath.Join(srcPath, "..", "extracted")
	}

	destOgPath := filepath.Join(outputDir, "og")

	err := shared.CopyCleanDir(srcPath, destOgPath)
	if err != nil {
		return fmt.Errorf("couldn't copy clean directory: %w", err)
	}

	err = qdecomp(filepath.Join(srcPath, "TEXTS.FIL"), filepath.Join(outputDir, "texts.txt"))
	if err != nil {
		return fmt.Errorf("couldn't decompile TEXTS.FIL: %w", err)
	}

	err = qdecomp(filepath.Join(srcPath, "RESOURCE.FIL"), filepath.Join(outputDir, "resource.txt"))
	if err != nil {
		return fmt.Errorf("couldn't decompile RESOURCE.FIL: %w", err)
	}

	err = qgetStrings(filepath.Join(srcPath, "GAME.EXE"), filepath.Join(outputDir, "game_exe.txt"), allStrings)
	if err != nil {
		return fmt.Errorf("couldn't get strings from GAME.EXE: %w", err)
	}

	err = qgetStrings(filepath.Join(srcPath, "INSTALL.EXE"), filepath.Join(outputDir, "install_exe.txt"), allStrings)
	if err != nil {
		return fmt.Errorf("couldn't get strings from INSTALL.EXE: %w", err)
	}

	return nil
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	debug := flag.Bool("v", false, "Enable verbose debug output")
	allStrings := flag.Bool("all-strings", false, "Extract all strings (non-conservative mode)")
	outputDir := flag.String("o", "", "Output directory (default: ../extracted relative to source)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("QADAM Extract Tool v%s\n", version)
		os.Exit(0)
	}

	// Get remaining arguments after flag parsing
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v <original source directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %v -version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %v -v <original source directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %v --all-strings <original source directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %v -o <output_dir> <original source directory>\n", os.Args[0])
		os.Exit(1)
	}

	// Set debug mode globally if requested
	debugMode = *debug
	if debugMode {
		fmt.Println("DEBUG: Verbose mode enabled")
	}

	if *allStrings {
		fmt.Println("INFO: All-strings mode enabled (non-conservative extraction)")
	}

	if *outputDir != "" {
		fmt.Printf("INFO: Output directory: %s\n", *outputDir)
	}

	err := extract(args[0], *outputDir, *allStrings)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		shared.PauseIfNeeded("Extraction failed! Press Enter to continue...")
		os.Exit(1)
	}

	// Pause if running from Explorer so the window doesn't close immediately
	shared.PauseIfNeeded("Extraction succeeded! Press Enter to continue...")
}
