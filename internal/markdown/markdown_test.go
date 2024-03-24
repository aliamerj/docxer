package markdown

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestCreateMarkdownDocx(t *testing.T) {
	// Setup temporary directory for test output
	dir, err := os.MkdirTemp("", "docxTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up after the test

	// Define Markdown input for the document
	markdown := "# Markdown Title\nThis is some sample text."

	outputFilePath, err := CreateMarkdownDocx(dir, markdown)
	if err != nil {
		t.Fatalf("CreateMarkdownDocx failed: %v", err)
	}

	// Ensure the file exists
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Fatalf("The DOCX file was not created.")
	}

	// Open the created .docx (ZIP) file
	file, err := os.Open(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to open created .docx file: %v", err)
	}
	defer file.Close()

	// Read the file size for zip.NewReader
	fi, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to obtain file info: %v", err)
	}

	// Use fileSize as the size parameter for zip.NewReader
	zipReader, err := zip.NewReader(file, fi.Size())
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

			// Here, we check for the presence of key Markdown-converted content in the DOCX output.
			// This assumes the conversion process transforms Markdown titles and text accordingly.
			if !strings.Contains(string(content), "Markdown Title") || !strings.Contains(string(content), "This is some sample text.") {
				t.Errorf("document.xml does not contain the correct converted Markdown content")
			}
			break // Stop after finding document.xml
		}
	}

	if !found {
		t.Errorf("document.xml was not found in the zip")
	}
}

func TestApplyStyle_Title(t *testing.T) {
	markdownInput := `#$# My Document Title`
	templateContent := `<w:document>
	<w:body>
		<w:p>{{TITLE}}</w:p>
		<w:p>{{SECTION}}</w:p>
	</w:body>
</w:document>`

	expectedOutput := strings.Replace(templateContent, "{{TITLE}}", "My Document Title", 1)
	expectedOutput = strings.Replace(expectedOutput, "{{SECTION}}", "", -1)

	output := applyStyle(templateContent, markdownInput)

	// Verify that the output matches expected output
	if output != expectedOutput {
		t.Errorf("Title not processed correctly. Expected:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
func TestApplyStyle_HeadingsWithStyle(t *testing.T) {
	templateContent := `<w:document>
	<w:body>
		{{SECTION}}
	</w:body>
</w:document>`

	// Pattern to extract the applied style from the output
	stylePattern := regexp.MustCompile(`<w:pStyle w:val="([^"]+)" />`)

	for level := 1; level <= 6; level++ {
		markdownInput := fmt.Sprintf("%s Heading Level %d", strings.Repeat("#", level), level)
		expectedStyle := fmt.Sprintf("Heading%d", level)

		// Apply the style to the template content using the markdownInput
		output := applyStyle(templateContent, markdownInput)

		// Extract the applied style from the output
		matches := stylePattern.FindStringSubmatch(output)
		if len(matches) < 2 {
			t.Fatalf("Failed to extract style for heading level %d from output", level)
		}
		appliedStyle := matches[1]

		// Verify that the applied style matches the expected style for the current heading level
		if appliedStyle != expectedStyle {
			t.Errorf("Incorrect style applied for heading level %d. Expected: %s, Got: %s", level, expectedStyle, appliedStyle)
		}
	}
}
