package document

import (
	"os"
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	// Scenario 1: Test empty path
	err := validateFilePath("")
	if err == nil {
		t.Errorf("validateFilePath(\"\") = %v, want error", err)
	}
	// Scenario 2: Test non-existing path
	nonExistPath := "path/that/does/not/exist"
	err = validateFilePath(nonExistPath)
	if err == nil {
		t.Errorf("validateFilePath(\"%s\") = %v, want error", nonExistPath, err)
	}

	// Scenario 3: Test with a valid directory
	// Create a temporary directory and defer its removal
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)

	err = validateFilePath(tempDir)
	if err != nil {
		t.Errorf("validateFilePath(\"%s\") = %v, want no error", tempDir, err)
	}

	// Scenario 4: Test with a file instead of a directory
	// Create a temporary file and defer its removal
	tempFile, err := os.CreateTemp(".", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)
	tempFile.Close()

	err = validateFilePath(tempFilePath)
	if err == nil {
		t.Errorf("validateFilePath(\"%s\") = %v, want error", tempFilePath, err)
	}
}
