package shared

import (
	"fmt"
	"os"
)

// IsRunningFromExplorer checks if the program is running from Windows Explorer
// or being redirected from a file (as opposed to running from command line)
func IsRunningFromExplorer() bool {
	// Check if stdin is a terminal - if not, we're likely running from Explorer
	// or being redirected from a file
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return true // Assume Explorer if we can't check
	}

	// If stdin is not a character device, we're likely running from Explorer
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

// PauseIfNeeded pauses execution if running from Explorer. If a message is provided, it is shown before the prompt.
func PauseIfNeeded(msg ...string) {
	if IsRunningFromExplorer() {
		if len(msg) > 0 && msg[0] != "" {
			fmt.Println(msg[0])
		}
		fmt.Println("Press Enter to continue...")
		fmt.Scanln() // Wait for Enter key
	}
}
