package markdown

import (
	"archive/zip"
	"io"
	"os"
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

func TestApplyStyle(t *testing.T) {
	tests := []struct {
		name            string
		fileContent     string
		markdownText    string
		expectedContent string
	}{
		{
			name:            "Single Title",
			fileContent:     "Document Title: {{TITLE}}\nContent:\n{{SECTION}}",
			markdownText:    "#$# My Title",
			expectedContent: "Document Title: My Title\nContent:\n",
		},
		{
			name:            "Title and Sections",
			fileContent:     "Document Title: {{TITLE}}\nContent:\n{{SECTION}}",
			markdownText:    "#$# My Title\n# Heading 1\nSome text.\n## Heading 2\nMore text.",
			expectedContent: "Document Title: My Title\nContent:\n" + `<w:p><w:pPr><w:pStyle w:val="Heading1"/></w:pPr><w:r><w:t xml:space="preserve">Heading 1</w:t></w:r></w:p><w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:t xml:space="preserve">Some text.</w:t></w:r></w:p><w:p><w:pPr><w:pStyle w:val="Heading2"/></w:pPr><w:r><w:t xml:space="preserve">Heading 2</w:t></w:r></w:p><w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:t xml:space="preserve">More text.</w:t></w:r></w:p>`,
		},

		{
			name:         "Normal Text",
			fileContent:  "{{TITLE}}\n{{SECTION}}",
			markdownText: "This is just some normal text without a title.",
			expectedContent: "\n" +
				`<w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:t xml:space="preserve">This is just some normal text without a title.</w:t></w:r></w:p>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := applyStyle(tc.fileContent, tc.markdownText)
			if got != tc.expectedContent {
				t.Errorf("TestApplyStyle %s failed:\nExpected:\n%s\nGot:\n%s", tc.name, tc.expectedContent, got)
			}
		})
	}
}

// TestDetermineStyle tests the determineStyle function for various input lines.
func TestDetermineStyle(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		line           string
		expectedStyle  string
		expectedMarker string
	}{
		{"Title", "#$# A title", "Title", "#$# "},
		{"Heading 1", "# Heading 1", "Heading1", "# "},
		{"Heading 2", "## Heading 2", "Heading2", "## "},
		{"Heading 3", "### Heading 3", "Heading3", "### "},
		{"Normal Text", "Just some normal text.", "Normal", ""},
	}

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style, placeholder := determineStyle(tt.line)
			if style != tt.expectedStyle || placeholder != tt.expectedMarker {
				t.Errorf("determineStyle(%q) got %q, %q; want %q, %q", tt.line, style, placeholder, tt.expectedStyle, tt.expectedMarker)
			}
		})
	}
}

func TestProcessTextFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No Styling",
			input:    "This text has no styling.",
			expected: `<w:r><w:t xml:space="preserve">This text has no styling.</w:t></w:r>`,
		},
		{
			name:     "Bold Styling",
			input:    "This text is **bold**.",
			expected: `<w:r><w:t xml:space="preserve">This </w:t></w:r> <w:r><w:t xml:space="preserve">text </w:t></w:r> <w:r><w:t xml:space="preserve">is </w:t></w:r> <w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">bold </w:t></w:r>.`,
		},
		{
			name:     "Italic Styling",
			input:    "This text is *italic*.",
			expected: `<w:r><w:t xml:space="preserve">This </w:t></w:r> <w:r><w:t xml:space="preserve">text </w:t></w:r> <w:r><w:t xml:space="preserve">is </w:t></w:r> <w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">italic </w:t></w:r>.`,
		},
		{
			name:     "Bold and Italic Styling",
			input:    "This is ***bold and italic*** text.",
			expected: `<w:r><w:t xml:space="preserve">This </w:t></w:r> <w:r><w:t xml:space="preserve">is </w:t></w:r> <w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">* </w:t></w:r>bold <w:r><w:t xml:space="preserve">and </w:t></w:r> italic<w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">* </w:t></w:r> <w:r><w:t xml:space="preserve">text. </w:t></w:r>`,
		},
		{
			name:     "Mixed Styling",
			input:    "This is **bold**, *italic*, and ***both***.",
			expected: `<w:r><w:t xml:space="preserve">This </w:t></w:r> <w:r><w:t xml:space="preserve">is </w:t></w:r> <w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">bold </w:t></w:r>, <w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">italic </w:t></w:r>, <w:r><w:t xml:space="preserve">and </w:t></w:r> <w:r><w:rPr><w:b/><w:i/></w:rPr><w:t xml:space="preserve">both </w:t></w:r>.`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := processTextFormatting(tc.input)
			if got != tc.expected {
				t.Errorf("Failed %s: Expected %s, got %s", tc.name, tc.expected, got)
			}
		})
	}
}

func TestNewSection(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		body     string
		expected string
	}{
		{
			name:     "Normal Text",
			style:    "Normal",
			body:     "Just some normal text.",
			expected: `<w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:t xml:space="preserve">Just some normal text.</w:t></w:r></w:p>`,
		},
		{
			name:     "Bold Text",
			style:    "Normal",
			body:     "**Bold** text.",
			expected: `<w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">Bold </w:t></w:r> <w:r><w:t xml:space="preserve">text. </w:t></w:r></w:p>`,
		},
		{
			name:     "Italic Text",
			style:    "Normal",
			body:     "*Italic* text.",
			expected: `<w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">Italic </w:t></w:r> <w:r><w:t xml:space="preserve">text. </w:t></w:r></w:p>`,
		},
		{
			name:     "Bold and Italic Text",
			style:    "Normal",
			body:     "***Bold and Italic*** text.",
			expected: `<w:p><w:pPr><w:pStyle w:val="Normal"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">* </w:t></w:r>Bold <w:r><w:t xml:space="preserve">and </w:t></w:r> Italic<w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">* </w:t></w:r> <w:r><w:t xml:space="preserve">text. </w:t></w:r></w:p>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := newSection(tc.style, tc.body)
			if got != tc.expected {
				t.Errorf("TestNewSection %s failed: Expected %s, got %s", tc.name, tc.expected, got)
			}
		})
	}
}
func TestEscapeXMLChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"This & That", "This &amp; That"},
		{"<Tag>", "&lt;Tag&gt;"},
		{"'Single' & \"Double\"", "&apos;Single&apos; &amp; &quot;Double&quot;"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := escapeXMLChars(test.input)
			if got != test.expected {
				t.Errorf("escapeXMLChars(%q) = %q, want %q", test.input, got, test.expected)
			}
		})
	}
}
