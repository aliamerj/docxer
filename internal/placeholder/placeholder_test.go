package placeholder

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"testing"
)

func TestUpdateDocx(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "testUpdateDocx")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Prepare a test DOCX file with placeholders
	testFilePath := tempDir + "/test.docx"
	createTestDocx(testFilePath, "Original content with {{TITLE}} and {{BODY}}")

	// Define replacements
	replacements := map[string]string{
		"TITLE": "Updated Title",
		"BODY":  "Updated Body",
	}

	// Run the UpdateDocx function
	err = UpdateDocx(testFilePath, replacements)
	if err != nil {
		t.Errorf("UpdateDocx returned an error: %v", err)
	}

	// Verify file content
	verifyUpdatedDocx(t, testFilePath, "Original content with Updated Title and Updated Body")
}

// createTestDocx creates a simple DOCX file with the specified content
func createTestDocx(filePath string, content string) {
	// Create a file
	file, err := os.Create(filePath)
	if err != nil {
		panic(err) // Panic in helper function as it's setup, not test logic
	}
	defer file.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Create a new file in the zip archive
	writer, err := zipWriter.Create("word/document.xml")
	if err != nil {
		panic(err)
	}

	// Write content to the new file
	_, err = writer.Write([]byte(content))
	if err != nil {
		panic(err)
	}
}

// verifyUpdatedDocx checks the content of the updated DOCX file
func verifyUpdatedDocx(t *testing.T, filePath string, expectedContent string) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		t.Fatalf("Failed to open updated DOCX file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("Failed to open document.xml: %v", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("Failed to read content from document.xml: %v", err)
			}

			if strings.TrimSpace(string(content)) != expectedContent {
				t.Errorf("Content did not match expected. Got: %s, Want: %s", content, expectedContent)
			}
			return
		}
	}
	t.Errorf("document.xml not found in ZIP")
}
