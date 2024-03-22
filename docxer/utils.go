package docxer

import (
	"fmt"
	"os"
)

func validateFilePath(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if the path exists and is a directory
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the specified directory does not exist: %s", dirPath)
		}
		return fmt.Errorf("error checking the specified directory: %v", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("the specified path is not a directory: %s", dirPath)
	}

	return nil
}
