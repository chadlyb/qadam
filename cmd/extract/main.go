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

func extract(srcPath string, allStrings bool) error {
	destPath := filepath.Join(srcPath, "..", "extracted")
	destOgPath := filepath.Join(destPath, "og")

	err := shared.CopyCleanDir(srcPath, destOgPath)
	if err != nil {
		return fmt.Errorf("couldn't copy clean directory: %w", err)
	}

	err = qdecomp(filepath.Join(srcPath, "TEXTS.FIL"), filepath.Join(destPath, "texts.txt"))
	if err != nil {
		return fmt.Errorf("couldn't decompile TEXTS.FIL: %w", err)
	}

	err = qdecomp(filepath.Join(srcPath, "RESOURCE.FIL"), filepath.Join(destPath, "resource.txt"))
	if err != nil {
		return fmt.Errorf("couldn't decompile RESOURCE.FIL: %w", err)
	}

	err = qgetStrings(filepath.Join(srcPath, "GAME.EXE"), filepath.Join(destPath, "game_exe.txt"), allStrings)
	if err != nil {
		return fmt.Errorf("couldn't get strings from GAME.EXE: %w", err)
	}

	err = qgetStrings(filepath.Join(srcPath, "INSTALL.EXE"), filepath.Join(destPath, "install_exe.txt"), allStrings)
	if err != nil {
		return fmt.Errorf("couldn't get strings from INSTALL.EXE: %w", err)
	}

	return nil
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	debug := flag.Bool("v", false, "Enable verbose debug output")
	allStrings := flag.Bool("all-strings", false, "Extract all strings (non-conservative mode)")
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

	err := extract(args[0], *allStrings)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		shared.PauseIfNeeded("Extraction failed! Press Enter to continue...")
		os.Exit(1)
	}

	// Pause if running from Explorer so the window doesn't close immediately
	shared.PauseIfNeeded("Extraction succeeded! Press Enter to continue...")
}
