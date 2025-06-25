package shared

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copies all files in srcPath to the directory tgtPath
// Directory is created if it does not exist
// ALL FILES IN TGTPATH ARE DESTROYED if they do exist!
func CopyCleanDir(srcPath string, tgtPath string) error {
	// Ensure source exists
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("cannot stat srcPath: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("srcPath is not a directory")
	}

	// If tgtPath exists, delete all contents
	if _, err := os.Stat(tgtPath); err == nil {
		err = os.RemoveAll(tgtPath)
		if err != nil {
			return fmt.Errorf("failed to clear tgtPath: %w", err)
		}
	}

	// Create tgtPath
	if err := os.MkdirAll(tgtPath, 0755); err != nil {
		return fmt.Errorf("failed to create tgtPath: %w", err)
	}

	// Walk source and copy files
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}

		// Construct target path
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return fmt.Errorf("rel path error: %w", err)
		}
		target := filepath.Join(tgtPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		// Copy file
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("open dst: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("copy contents: %w", err)
	}

	return nil
}
