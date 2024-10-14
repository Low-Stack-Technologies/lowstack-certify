package provider_utils

import (
	"fmt"
	"os"
)

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0600); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Check if directory is writable
	if err := os.Chmod(path, 0700); err != nil {
		return fmt.Errorf("failed to set directory permissions: %w", err)
	}

	return nil
}
