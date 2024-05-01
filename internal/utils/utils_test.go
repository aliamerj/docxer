package utils

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"
	"testing/fstest"
)

func TestValidateFilePath(t *testing.T) {
	// Scenario 1: Test empty path
	err := ValidateFilePath("")
	if err == nil {
		t.Errorf("validateFilePath(\"\") = %v, want error", err)
	}
	// Scenario 2: Test non-existing path
	nonExistPath := "path/that/does/not/exist"
	err = ValidateFilePath(nonExistPath)
	if err == nil {
		t.Errorf("validateFilePath(\"%s\") = %v, want error", nonExistPath, err)
	}

	// Scenario 3: Test with a valid directory
	// Create a temporary directory and defer its removal
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)

	err = ValidateFilePath(tempDir)
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

	err = ValidateFilePath(tempFilePath)
	if err == nil {
		t.Errorf("validateFilePath(\"%s\") = %v, want error", tempFilePath, err)
	}
}

// MockDocxWriter simply appends "Processed" to any input text.
func MockDocxWriter(input string) string {
	return input + "Processed"
}

func TestCreateDocx(t *testing.T) {
	// Create an in-memory file system (mock FS) with template files.
	memFS := fstest.MapFS{
		"template/document.xml": &fstest.MapFile{Data: []byte("Content of document.xml")},
		// Add more files as needed.
	}

	// Create an in-memory ZIP writer to capture the output.
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Call the function under test.
	err := CreateDocx(memFS, zipWriter, MockDocxWriter)
	if err != nil {
		t.Fatalf("CreateDocx failed: %v", err)
	}

	// Ensure everything is written to the buffer.
	if err := zipWriter.Close(); err != nil {
		t.Fatalf("Failed to close zip writer: %v", err)
	}

	// Read back the ZIP content to verify.
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("Failed to read zip content: %v", err)
	}

	// Verify the contents of the document.xml in the ZIP file.
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("Failed to open %s: %v", f.Name, err)
			}
			var fileContent bytes.Buffer
			if _, err := fileContent.ReadFrom(rc); err != nil {
				t.Fatalf("Failed to read %s: %v", f.Name, err)
			}
			rc.Close() // Close the file after reading

			expectedContent := "Content of document.xmlProcessed"
			if fileContent.String() != expectedContent {
				t.Errorf("Expected content %q; got %q", expectedContent, fileContent.String())
			}
		}
	}
}
