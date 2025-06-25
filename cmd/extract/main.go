package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chadlyb/qadam/shared"
)

func extract(srcPath string) error {
	destPath := filepath.Join(srcPath, "..", "extracted")
	destOgPath := filepath.Join(destPath, "og")

	err := shared.CopyCleanDir(srcPath, destOgPath)
	if err != nil {
		return err
	}

	err = qdecomp(filepath.Join(srcPath, "TEXTS.FIL"), filepath.Join(destPath, "texts.txt"))
	if err != nil {
		return err
	}

	err = qdecomp(filepath.Join(srcPath, "RESOURCE.FIL"), filepath.Join(destPath, "resource.txt"))
	if err != nil {
		return err
	}

	err = qgetStrings(filepath.Join(srcPath, "GAME.EXE"), filepath.Join(destPath, "game_exe.txt"))
	if err != nil {
		return err
	}

	err = qgetStrings(filepath.Join(srcPath, "INSTALL.EXE"), filepath.Join(destPath, "install_exe.txt"))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <original source directory>\n", os.Args[0])
		os.Exit(1)
	}

	err := extract(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
