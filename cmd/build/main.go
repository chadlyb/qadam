package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chadlyb/qadam/shared"
)

func build(srcPath string) error {
	srcOgPath := filepath.Join(srcPath, "og")
	destPath := filepath.Join(srcPath, "..", "built")

	err := shared.CopyCleanDir(srcOgPath, destPath)
	if err != nil {
		return err
	}

	textsFil := filepath.Join(destPath, "TEXTS.FIL")
	resourceFil := filepath.Join(destPath, "RESOURCE.FIL")
	gameExe := filepath.Join(destPath, "GAME.EXE")

	err = qcompile(filepath.Join(srcPath, "texts.txt"), textsFil)
	if err != nil {
		return err
	}

	err = qcompile(filepath.Join(srcPath, "resource.txt"), resourceFil)
	if err != nil {
		return err
	}

	err = qpatchStrings(filepath.Join(srcOgPath, "GAME.EXE"), gameExe, filepath.Join(srcPath, "game_exe.txt"))
	if err != nil {
		return err
	}

	err = qpatchStrings(filepath.Join(srcOgPath, "INSTALL.EXE"), filepath.Join(destPath, "INSTALL.EXE"), filepath.Join(srcPath, "install_exe.txt"))
	if err != nil {
		return err
	}

	err = fixgame(gameExe, gameExe, textsFil, resourceFil)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <extracted directory>\n", os.Args[0])
		os.Exit(1)
	}

	err := build(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
