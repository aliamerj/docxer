package placeholder

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"testing"
)

func TestTextPlaceholderWriter_BasicReplacement(t *testing.T) {
	// Setup test input and output
	inputContent := "Hello {{NAME}}, welcome to {{PLACE}}."
	replacements := map[string]string{
		"NAME":  "John Doe",
		"PLACE": "GoLand",
	}
	expectedOutput := "Hello John Doe, welcome to GoLand."

	// Create the placeholder writer action
	action := TextPlaceholderWriter(replacements)
	docxWriter := action()

	// Execute the writer
	outputContent := docxWriter(inputContent)

	// Verify the output
	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestTextPlaceholderWriter_NoReplacements(t *testing.T) {
	// Setup test input and output
	inputContent := "Hello there."
	replacements := map[string]string{
		"UNUSED": "Unused",
	}
	expectedOutput := "Hello there."

	// Create the placeholder writer action
	action := TextPlaceholderWriter(replacements)
	docxWriter := action()

	// Execute the writer
	outputContent := docxWriter(inputContent)

	// Verify the output
	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestTextPlaceholderWriter_MultipleOccurrences(t *testing.T) {
	// Setup test input and output
	inputContent := "Hello {{NAME}}, you are {{NAME}} right?"
	replacements := map[string]string{
		"NAME": "Alice",
	}
	expectedOutput := "Hello Alice, you are Alice right?"

	// Create the placeholder writer action
	action := TextPlaceholderWriter(replacements)
	docxWriter := action()

	// Execute the writer
	outputContent := docxWriter(inputContent)

	// Verify the output
	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestTextPlaceholderWriter_EmptyStrings(t *testing.T) {
	// Setup test input and output
	inputContent := "Hi {{NAME}}!"
	replacements := map[string]string{
		"NAME": "",
	}
	expectedOutput := "Hi !"

	// Create the placeholder writer action
	action := TextPlaceholderWriter(replacements)
	docxWriter := action()

	// Execute the writer
	outputContent := docxWriter(inputContent)

	// Verify the output
	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestLoopPlaceholderWriter_BasicLoop(t *testing.T) {
	// Setup test input and output
	inputContent := "Items:\n{{#each items}}- {{NAME}}, ${{PRICE}}\n{{/each}}"
	data := map[string]interface{}{
		"items": []map[string]string{
			{"NAME": "Item 1", "PRICE": "10"},
			{"NAME": "Item 2", "PRICE": "20"},
		},
	}
	expectedOutput := "Items:\n- Item 1, $10\n- Item 2, $20\n"

	// Create the placeholder writer action
	action := LoopPlaceholderWriter(data)
	docxWriter := action()

	// Execute the writer
	outputContent := docxWriter(inputContent)

	// Verify the output
	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestLoopPlaceholderWriter_MultipleLoops(t *testing.T) {
	inputContent := "Items:\n{{#each items}}- {{NAME}}, ${{PRICE}}\n{{/each}}Services:\n{{#each services}}- {{SERVICE}}: ${{COST}}\n{{/each}}"
	data := map[string]interface{}{
		"items": []map[string]string{
			{"NAME": "Item 1", "PRICE": "10"},
			{"NAME": "Item 2", "PRICE": "20"},
		},
		"services": []map[string]string{
			{"SERVICE": "Delivery", "COST": "5"},
			{"SERVICE": "Gift Wrap", "COST": "2"},
		},
	}
	expectedOutput := "Items:\n- Item 1, $10\n- Item 2, $20\nServices:\n- Delivery: $5\n- Gift Wrap: $2\n"

	action := LoopPlaceholderWriter(data)
	docxWriter := action()

	outputContent := docxWriter(inputContent)

	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestLoopPlaceholderWriter_MissingMarkers(t *testing.T) {
	inputContent := "There are no loop markers here."
	data := map[string]interface{}{
		"items": []map[string]string{
			{"NAME": "Item 1", "PRICE": "10"},
		},
	}
	expectedOutput := "There are no loop markers here."

	action := LoopPlaceholderWriter(data)
	docxWriter := action()

	outputContent := docxWriter(inputContent)

	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestLoopPlaceholderWriter_IncorrectDataFormat(t *testing.T) {
	inputContent := "Items:\n{{#each items}}- {{NAME}}, ${{PRICE}}\n{{/each}}"
	data := map[string]interface{}{
		"items": "Not a slice of maps",
	}
	expectedOutput := "Items:\n{{#each items}}- {{NAME}}, ${{PRICE}}\n{{/each}}"

	action := LoopPlaceholderWriter(data)
	docxWriter := action()

	outputContent := docxWriter(inputContent)

	if outputContent != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, outputContent)
	}
}

func TestUpdateDocx_Basic(t *testing.T) {
	// Setup a temporary directory
	tempDir, err := os.MkdirTemp("", "testUpdateDocx_Basic")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test DOCX with placeholders
	testFilePath := tempDir + "/test.docx"
	createTestDocx(testFilePath, "Original content with {{TITLE}} and {{BODY}}")

	// Create an action to replace placeholders
	action := TextPlaceholderWriter(map[string]string{
		"TITLE": "Updated Title",
		"BODY":  "Updated Body",
	})

	// Call UpdateDocx with the action
	if err := UpdateDocx(testFilePath, action); err != nil {
		t.Errorf("UpdateDocx returned an error: %v", err)
	}

	// Verify the content
	verifyUpdatedDocx(t, testFilePath, "Original content with Updated Title and Updated Body")
}

func TestUpdateDocx_MultipleReplacements(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testUpdateDocx_MultipleReplacements")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := tempDir + "/test.docx"
	createTestDocx(testFilePath, "Info: {{NAME}}, Age: {{AGE}}")

	action := TextPlaceholderWriter(map[string]string{
		"NAME": "Alice",
		"AGE":  "30",
	})

	if err := UpdateDocx(testFilePath, action); err != nil {
		t.Errorf("UpdateDocx returned an error: %v", err)
	}

	expectedContent := "Info: Alice, Age: 30"
	verifyUpdatedDocx(t, testFilePath, expectedContent)

}

func TestUpdateDocx_NoReplacements(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testUpdateDocx_NoReplacements")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := tempDir + "/test.docx"
	createTestDocx(testFilePath, "Hello there.")

	action := TextPlaceholderWriter(map[string]string{
		"UNUSED": "Nothing",
	})

	if err := UpdateDocx(testFilePath, action); err != nil {
		t.Errorf("UpdateDocx returned an error: %v", err)
	}

	verifyUpdatedDocx(t, testFilePath, "Hello there.")
}

func TestUpdateDocx_FileDoesNotExist(t *testing.T) {
	nonexistentFile := "/path/to/nonexistent.docx"
	action := TextPlaceholderWriter(map[string]string{
		"TITLE": "Title",
	})

	err := UpdateDocx(nonexistentFile, action)
	if err == nil {
		t.Errorf("Expected error for non-existent file")
	}
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
