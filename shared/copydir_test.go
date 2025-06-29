package shared

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyCleanDir(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	t.Run("copy simple directory structure", func(t *testing.T) {
		// Create source directory structure
		err := os.MkdirAll(srcDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		// Create some test files
		testFiles := map[string]string{
			"file1.txt":        "Hello, World!",
			"file2.txt":        "Test content",
			"subdir/file3.txt": "Nested file content",
			"subdir/file4.txt": "Another nested file",
		}

		for filePath, content := range testFiles {
			fullPath := filepath.Join(srcDir, filePath)
			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			if err != nil {
				t.Fatalf("Failed to create subdirectory: %v", err)
			}

			err = os.WriteFile(fullPath, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file %s: %v", filePath, err)
			}
		}

		// Copy directory
		err = CopyCleanDir(srcDir, dstDir)
		if err != nil {
			t.Fatalf("CopyCleanDir failed: %v", err)
		}

		// Verify all files were copied correctly
		for filePath, expectedContent := range testFiles {
			fullPath := filepath.Join(dstDir, filePath)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read copied file %s: %v", filePath, err)
				continue
			}

			if string(content) != expectedContent {
				t.Errorf("File %s content mismatch. Expected: %q, Got: %q", filePath, expectedContent, string(content))
			}
		}

		// Verify directory structure
		expectedDirs := []string{"", "subdir"}
		for _, dir := range expectedDirs {
			dirPath := filepath.Join(dstDir, dir)
			info, err := os.Stat(dirPath)
			if err != nil {
				t.Errorf("Failed to stat directory %s: %v", dirPath, err)
				continue
			}
			if !info.IsDir() {
				t.Errorf("Expected %s to be a directory", dirPath)
			}
		}
	})

	t.Run("copy to existing directory (should clean first)", func(t *testing.T) {
		// Create a file in the destination directory
		existingFile := filepath.Join(dstDir, "existing.txt")
		err := os.WriteFile(existingFile, []byte("This should be deleted"), 0644)
		if err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		// Create a new source with different content
		newSrcDir := filepath.Join(tempDir, "src2")
		err = os.MkdirAll(newSrcDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create new source directory: %v", err)
		}

		newFile := filepath.Join(newSrcDir, "newfile.txt")
		err = os.WriteFile(newFile, []byte("New content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create new test file: %v", err)
		}

		// Copy to existing destination
		err = CopyCleanDir(newSrcDir, dstDir)
		if err != nil {
			t.Fatalf("CopyCleanDir failed: %v", err)
		}

		// Verify old file was deleted
		if _, err := os.Stat(existingFile); err == nil {
			t.Error("Existing file was not deleted")
		}

		// Verify new file was copied
		newFileDst := filepath.Join(dstDir, "newfile.txt")
		content, err := os.ReadFile(newFileDst)
		if err != nil {
			t.Errorf("Failed to read new file: %v", err)
		} else if string(content) != "New content" {
			t.Errorf("New file content mismatch. Expected: %q, Got: %q", "New content", string(content))
		}
	})

	t.Run("copy empty directory", func(t *testing.T) {
		emptySrcDir := filepath.Join(tempDir, "empty_src")
		emptyDstDir := filepath.Join(tempDir, "empty_dst")

		err := os.MkdirAll(emptySrcDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create empty source directory: %v", err)
		}

		err = CopyCleanDir(emptySrcDir, emptyDstDir)
		if err != nil {
			t.Fatalf("CopyCleanDir failed: %v", err)
		}

		// Verify destination directory was created
		info, err := os.Stat(emptyDstDir)
		if err != nil {
			t.Errorf("Failed to stat destination directory: %v", err)
		} else if !info.IsDir() {
			t.Error("Destination is not a directory")
		}
	})

	t.Run("copy directory with subdirectories", func(t *testing.T) {
		nestedSrcDir := filepath.Join(tempDir, "nested_src")
		nestedDstDir := filepath.Join(tempDir, "nested_dst")

		// Create nested directory structure
		nestedFiles := map[string]string{
			"level1/file1.txt":               "Level 1 file",
			"level1/level2/file2.txt":        "Level 2 file",
			"level1/level2/level3/file3.txt": "Level 3 file",
		}

		for filePath, content := range nestedFiles {
			fullPath := filepath.Join(nestedSrcDir, filePath)
			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			if err != nil {
				t.Fatalf("Failed to create nested directory: %v", err)
			}

			err = os.WriteFile(fullPath, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create nested file %s: %v", filePath, err)
			}
		}

		err := CopyCleanDir(nestedSrcDir, nestedDstDir)
		if err != nil {
			t.Fatalf("CopyCleanDir failed: %v", err)
		}

		// Verify nested structure was copied
		for filePath, expectedContent := range nestedFiles {
			fullPath := filepath.Join(nestedDstDir, filePath)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read nested file %s: %v", filePath, err)
				continue
			}

			if string(content) != expectedContent {
				t.Errorf("Nested file %s content mismatch. Expected: %q, Got: %q", filePath, expectedContent, string(content))
			}
		}
	})
}

func TestCopyCleanDirErrors(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("source does not exist", func(t *testing.T) {
		nonexistentSrc := filepath.Join(tempDir, "nonexistent")
		dstDir := filepath.Join(tempDir, "dst")

		err := CopyCleanDir(nonexistentSrc, dstDir)
		if err == nil {
			t.Error("Expected error when source does not exist")
		}
		if !strings.Contains(err.Error(), "cannot stat srcPath") {
			t.Errorf("Expected error about stat srcPath, got: %v", err)
		}
	})

	t.Run("source is not a directory", func(t *testing.T) {
		// Create a file instead of directory
		srcFile := filepath.Join(tempDir, "srcfile.txt")
		err := os.WriteFile(srcFile, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		dstDir := filepath.Join(tempDir, "dst")

		err = CopyCleanDir(srcFile, dstDir)
		if err == nil {
			t.Error("Expected error when source is not a directory")
		}
		if !strings.Contains(err.Error(), "srcPath is not a directory") {
			t.Errorf("Expected error about srcPath not being a directory, got: %v", err)
		}
	})

	t.Run("cannot create destination directory", func(t *testing.T) {
		srcDir := filepath.Join(tempDir, "src")
		err := os.MkdirAll(srcDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		// Try to copy to a path that has a file in the middle of the path
		// This should fail because we can't create a directory where a file exists
		intermediateFile := filepath.Join(tempDir, "intermediate")
		err = os.WriteFile(intermediateFile, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create intermediate file: %v", err)
		}

		// Try to copy to a path that goes through the file
		dstPath := filepath.Join(intermediateFile, "dst")
		err = CopyCleanDir(srcDir, dstPath)
		if err == nil {
			t.Error("Expected error when cannot create destination directory")
		} else if !strings.Contains(err.Error(), "failed to create tgtPath") {
			t.Errorf("Expected error about creating tgtPath, got: %v", err)
		}
	})
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("copy simple file", func(t *testing.T) {
		srcFile := filepath.Join(tempDir, "src.txt")
		dstFile := filepath.Join(tempDir, "dst.txt")

		// Create source file
		content := "Hello, World! This is a test file."
		err := os.WriteFile(srcFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Get file info for permissions
		info, err := os.Stat(srcFile)
		if err != nil {
			t.Fatalf("Failed to stat source file: %v", err)
		}

		// Copy file
		err = copyFile(srcFile, dstFile, info.Mode())
		if err != nil {
			t.Fatalf("copyFile failed: %v", err)
		}

		// Verify file was copied
		copiedContent, err := os.ReadFile(dstFile)
		if err != nil {
			t.Errorf("Failed to read copied file: %v", err)
		} else if string(copiedContent) != content {
			t.Errorf("File content mismatch. Expected: %q, Got: %q", content, string(copiedContent))
		}

		// Verify permissions were preserved
		dstInfo, err := os.Stat(dstFile)
		if err != nil {
			t.Errorf("Failed to stat destination file: %v", err)
		} else if dstInfo.Mode() != info.Mode() {
			t.Errorf("Permission mismatch. Expected: %v, Got: %v", info.Mode(), dstInfo.Mode())
		}
	})

	t.Run("copy large file", func(t *testing.T) {
		srcFile := filepath.Join(tempDir, "large_src.txt")
		dstFile := filepath.Join(tempDir, "large_dst.txt")

		// Create a larger file (1MB)
		largeContent := strings.Repeat("0123456789", 100000) // 1MB of data
		err := os.WriteFile(srcFile, []byte(largeContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create large source file: %v", err)
		}

		info, err := os.Stat(srcFile)
		if err != nil {
			t.Fatalf("Failed to stat large source file: %v", err)
		}

		err = copyFile(srcFile, dstFile, info.Mode())
		if err != nil {
			t.Fatalf("copyFile failed for large file: %v", err)
		}

		// Verify large file was copied correctly
		copiedContent, err := os.ReadFile(dstFile)
		if err != nil {
			t.Errorf("Failed to read large copied file: %v", err)
		} else if string(copiedContent) != largeContent {
			t.Error("Large file content mismatch")
		}
	})

	t.Run("copy file with special permissions", func(t *testing.T) {
		srcFile := filepath.Join(tempDir, "perm_src.txt")
		dstFile := filepath.Join(tempDir, "perm_dst.txt")

		// Create source file with executable permissions
		err := os.WriteFile(srcFile, []byte("executable file"), 0755)
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		info, err := os.Stat(srcFile)
		if err != nil {
			t.Fatalf("Failed to stat source file: %v", err)
		}

		err = copyFile(srcFile, dstFile, info.Mode())
		if err != nil {
			t.Fatalf("copyFile failed: %v", err)
		}

		// Verify permissions were preserved
		dstInfo, err := os.Stat(dstFile)
		if err != nil {
			t.Errorf("Failed to stat destination file: %v", err)
		} else if dstInfo.Mode() != info.Mode() {
			t.Errorf("Permission mismatch. Expected: %v, Got: %v", info.Mode(), dstInfo.Mode())
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		srcFile := filepath.Join(tempDir, "overwrite_src.txt")
		dstFile := filepath.Join(tempDir, "overwrite_dst.txt")

		// Create source file
		srcContent := "New content"
		err := os.WriteFile(srcFile, []byte(srcContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Create existing destination file
		oldContent := "Old content that should be overwritten"
		err = os.WriteFile(dstFile, []byte(oldContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create existing destination file: %v", err)
		}

		info, err := os.Stat(srcFile)
		if err != nil {
			t.Fatalf("Failed to stat source file: %v", err)
		}

		err = copyFile(srcFile, dstFile, info.Mode())
		if err != nil {
			t.Fatalf("copyFile failed: %v", err)
		}

		// Verify file was overwritten
		copiedContent, err := os.ReadFile(dstFile)
		if err != nil {
			t.Errorf("Failed to read overwritten file: %v", err)
		} else if string(copiedContent) != srcContent {
			t.Errorf("File content mismatch after overwrite. Expected: %q, Got: %q", srcContent, string(copiedContent))
		}
	})
}

func TestCopyFileErrors(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("source file does not exist", func(t *testing.T) {
		nonexistentSrc := filepath.Join(tempDir, "nonexistent.txt")
		dstFile := filepath.Join(tempDir, "dst.txt")

		err := copyFile(nonexistentSrc, dstFile, 0644)
		if err == nil {
			t.Error("Expected error when source file does not exist")
		}
		if !strings.Contains(err.Error(), "open src") {
			t.Errorf("Expected error about opening src, got: %v", err)
		}
	})

	t.Run("destination directory does not exist", func(t *testing.T) {
		srcFile := filepath.Join(tempDir, "src.txt")
		dstFile := filepath.Join(tempDir, "nonexistent", "dst.txt")

		// Create source file
		err := os.WriteFile(srcFile, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		err = copyFile(srcFile, dstFile, 0644)
		if err == nil {
			t.Error("Expected error when destination directory does not exist")
		}
		if !strings.Contains(err.Error(), "open dst") {
			t.Errorf("Expected error about opening dst, got: %v", err)
		}
	})
}

func TestCopyCleanDirWithSymlinks(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("copy directory with symlinks", func(t *testing.T) {
		srcDir := filepath.Join(tempDir, "src")
		dstDir := filepath.Join(tempDir, "dst")

		// Create source directory structure
		err := os.MkdirAll(srcDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		// Create a regular file
		regularFile := filepath.Join(srcDir, "regular.txt")
		err = os.WriteFile(regularFile, []byte("Regular file content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create regular file: %v", err)
		}

		// Create a symlink (if supported)
		symlinkFile := filepath.Join(srcDir, "symlink.txt")
		err = os.Symlink(regularFile, symlinkFile)
		if err != nil {
			// Skip this test if symlinks are not supported
			t.Skip("Symlinks not supported on this platform")
		}

		// Copy directory
		err = CopyCleanDir(srcDir, dstDir)
		if err != nil {
			t.Fatalf("CopyCleanDir failed: %v", err)
		}

		// Verify regular file was copied
		regularDst := filepath.Join(dstDir, "regular.txt")
		content, err := os.ReadFile(regularDst)
		if err != nil {
			t.Errorf("Failed to read copied regular file: %v", err)
		} else if string(content) != "Regular file content" {
			t.Errorf("Regular file content mismatch")
		}

		// Verify symlink was copied (should be a regular file, not a symlink)
		symlinkDst := filepath.Join(dstDir, "symlink.txt")
		content, err = os.ReadFile(symlinkDst)
		if err != nil {
			t.Errorf("Failed to read copied symlink file: %v", err)
		} else if string(content) != "Regular file content" {
			t.Errorf("Symlink file content mismatch")
		}

		// Verify it's not a symlink in the destination
		info, err := os.Lstat(symlinkDst)
		if err != nil {
			t.Errorf("Failed to lstat copied symlink file: %v", err)
		} else if info.Mode()&os.ModeSymlink != 0 {
			t.Error("Copied file should not be a symlink")
		}
	})
}

func TestCopyCleanDirPreservesFileModes(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("preserve file modes", func(t *testing.T) {
		srcDir := filepath.Join(tempDir, "src")
		dstDir := filepath.Join(tempDir, "dst")

		// Create source directory
		err := os.MkdirAll(srcDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		// Create files with different permissions
		testFiles := map[string]os.FileMode{
			"readonly.txt":   0444,
			"executable.txt": 0755,
			"normal.txt":     0644,
		}

		for fileName, mode := range testFiles {
			filePath := filepath.Join(srcDir, fileName)
			err := os.WriteFile(filePath, []byte("test"), mode)
			if err != nil {
				t.Fatalf("Failed to create test file %s: %v", fileName, err)
			}
		}

		// Copy directory
		err = CopyCleanDir(srcDir, dstDir)
		if err != nil {
			t.Fatalf("CopyCleanDir failed: %v", err)
		}

		// Verify file modes were preserved
		for fileName, expectedMode := range testFiles {
			filePath := filepath.Join(dstDir, fileName)
			info, err := os.Stat(filePath)
			if err != nil {
				t.Errorf("Failed to stat copied file %s: %v", fileName, err)
				continue
			}

			// Compare only the permission bits
			actualMode := info.Mode() & 0777
			expectedMode = expectedMode & 0777
			if actualMode != expectedMode {
				t.Errorf("File mode mismatch for %s. Expected: %v, Got: %v", fileName, expectedMode, actualMode)
			}
		}
	})
}
