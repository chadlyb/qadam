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

	err = fixgame(gameExe, gameExe, textsFil, resourceFil)
	if err != nil {
		return fmt.Errorf("failed to fix game executable: %w", err)
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
