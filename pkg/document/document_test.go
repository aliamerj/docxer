package document

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCreateNewDocx(t *testing.T) {
	// Setup temporary directory for test output
	dir, err := os.MkdirTemp("", "docxTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	// Call the function under test
	title := "Test Title"
	body := "Test Body"
	outputFilePath, err := CreateNewDocx(dir, title, body)
	if err != nil {
		t.Fatalf("CreateNewDocx failed: %v", err)
	}

	// Open the created .docx (ZIP) file
	file, err := os.Open(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to open created .docx file: %v", err)
	}
	defer file.Close()

	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		t.Fatalf("Failed to seek file: %v", err)
	}

	// Use fileSize as the size parameter for zip.NewReader
	zipReader, err := zip.NewReader(file, fileSize)
	if err != nil {
		t.Fatalf("Failed to read zip file: %v", err)
	}

	// Look for the document.xml file within the ZIP and check its contents
	found := false
	for _, zipFile := range zipReader.File {
		if zipFile.Name == "word/document.xml" {
			found = true
			rc, err := zipFile.Open()
			if err != nil {
				t.Fatalf("Failed to open document.xml: %v", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("Failed to read document.xml: %v", err)
			}

			if !strings.Contains(string(content), title) || !strings.Contains(string(content), body) {
				t.Errorf("document.xml does not contain the correct title and/or body")
			}
		}
	}

	if !found {
		t.Errorf("document.xml was not found in the zip")
	}
}
