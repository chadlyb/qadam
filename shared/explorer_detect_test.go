package shared

import (
	"os"
	"testing"
)

func TestIsRunningFromExplorer(t *testing.T) {
	// This test will behave differently depending on how it's run
	// When run from command line, it should return false
	// When run from IDE or redirected input, it might return true

	result := IsRunningFromExplorer()

	// We can't make strong assertions about the result since it depends on the environment
	// But we can verify the function doesn't panic and returns a boolean
	t.Logf("IsRunningFromExplorer() returned: %v", result)

	// Test that the function handles errors gracefully
	// We can't easily test the error case without mocking, but the function should not panic
}

func TestPauseIfNeeded(t *testing.T) {
	// This test verifies that PauseIfNeeded doesn't panic
	// In a test environment, it should not actually pause

	// Capture stdout to verify the message is printed when appropriate
	// For now, just verify it doesn't panic
	PauseIfNeeded()

	// If we get here without panicking, the test passes
	t.Log("PauseIfNeeded() completed without panicking")
}

func TestExplorerDetectionWithStdin(t *testing.T) {
	// Test the underlying logic by checking stdin characteristics
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		t.Logf("Could not stat stdin: %v", err)
		return
	}

	isCharDevice := (fileInfo.Mode() & os.ModeCharDevice) != 0
	t.Logf("Stdin is character device: %v", isCharDevice)
	t.Logf("Stdin mode: %v", fileInfo.Mode())
}
