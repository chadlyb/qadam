package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chadlyb/qadam/shared"
)

// Version will be set by the linker during build
var version = "dev"

func patchFileSizes(gameExePath, textsFilPath, resourceFilPath string) error {
	gameExeData, err := os.ReadFile(gameExePath)
	if err != nil {
		return fmt.Errorf("failed to read game executable: %w", err)
	}

	// Get file sizes for patching
	textsInfo, err := os.Lstat(textsFilPath)
	if err != nil {
		return fmt.Errorf("failed to get TEXTS.FIL size: %w", err)
	}
	resourceInfo, err := os.Lstat(resourceFilPath)
	if err != nil {
		return fmt.Errorf("failed to get RESOURCE.FIL size: %w", err)
	}

	// Patch the file sizes directly in the buffer
	const textsFilOffset = 0x0001A706
	const resourceFilOffset = 0x0001A6E6

	// Check bounds
	if textsFilOffset+4 > len(gameExeData) {
		return fmt.Errorf("TEXTS.FIL offset 0x%X out of bounds (size %d bytes)", textsFilOffset, len(gameExeData))
	}
	if resourceFilOffset+4 > len(gameExeData) {
		return fmt.Errorf("RESOURCE.FIL offset 0x%X out of bounds (size %d bytes)", resourceFilOffset, len(gameExeData))
	}

	// Patch the values in-place
	binary.LittleEndian.PutUint32(gameExeData[textsFilOffset:], uint32(textsInfo.Size()))
	binary.LittleEndian.PutUint32(gameExeData[resourceFilOffset:], uint32(resourceInfo.Size()))

	// Write the patched data back to the file
	err = os.WriteFile(gameExePath, gameExeData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write patched game executable: %w", err)
	}

	return nil
}

func build(srcPath string) error {
	srcOgPath := filepath.Join(srcPath, "og")
	destPath := filepath.Join(srcPath, "..", "built")

	err := shared.CopyCleanDir(srcOgPath, destPath)
	if err != nil {
		return fmt.Errorf("failed to copy clean directory: %w", err)
	}

	textsFil := filepath.Join(destPath, "TEXTS.FIL")
	resourceFil := filepath.Join(destPath, "RESOURCE.FIL")
	gameExe := filepath.Join(destPath, "GAME.EXE")

	err = qcompile(filepath.Join(srcPath, "texts.txt"), textsFil)
	if err != nil {
		return fmt.Errorf("failed to compile texts.txt: %w", err)
	}

	err = qcompile(filepath.Join(srcPath, "resource.txt"), resourceFil)
	if err != nil {
		return fmt.Errorf("failed to compile resource.txt: %w", err)
	}

	err = qpatchStrings(filepath.Join(srcOgPath, "GAME.EXE"), gameExe, filepath.Join(srcPath, "game_exe.txt"))
	if err != nil {
		return fmt.Errorf("failed to patch strings in GAME.EXE: %w", err)
	}

	err = qpatchStrings(filepath.Join(srcOgPath, "INSTALL.EXE"), filepath.Join(destPath, "INSTALL.EXE"), filepath.Join(srcPath, "install_exe.txt"))
	if err != nil {
		return fmt.Errorf("failed to patch strings in INSTALL.EXE: %w", err)
	}

	// Patch the game executable to have correct file sizes
	err = patchFileSizes(gameExe, textsFil, resourceFil)
	if err != nil {
		return fmt.Errorf("failed to patch file sizes: %w", err)
	}

	return nil
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("QADAM Build Tool v%s\n", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <extracted directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %v -version\n", os.Args[0])
		os.Exit(1)
	}

	err := build(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
