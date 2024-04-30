package placeholder

import (
	"archive/zip"
	"io"
	"os"
	"strings"

	"github.com/aliamerj/docxer/internal/utils"
)

type PlaceholderAction func() utils.DocxWriter

func UpdateDocx(filePath string, action PlaceholderAction) error {
	// Open the existing DOCX file for reading
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Create a temporary output file
	tempFilePath := filePath + ".tmp"
	outputFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	// Setup the docxWriter function for updating placeholders
	docxer := action() // docxPlaceholderWriter(replacements)

	// Process each file in the zip archive
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		fileContent, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		newFile, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}
		// Apply transformation if it's the document XML or other target files
		updatedContent := docxer(string(fileContent))
		if _, err = newFile.Write([]byte(updatedContent)); err != nil {
			return err
		}
	}

	// Replace the original file with the updated version
	if err := zipWriter.Close(); err != nil {
		return err
	}
	if err := outputFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempFilePath, filePath); err != nil {
		return err
	}

	return nil
}

// DocxPlaceholderWriter creates a function to replace placeholders with their corresponding replacements.
func TextPlaceholderWriter(replacements map[string]string) PlaceholderAction {
	return func() utils.DocxWriter {
		return func(fileContent string) string {
			updatedTemplate := fileContent
			for placeholder, replacement := range replacements {
				updatedTemplate = strings.ReplaceAll(updatedTemplate, "{{"+placeholder+"}}", replacement)
			}
			return updatedTemplate
		}
	}
}

func LoopPlaceholderWriter(data map[string]interface{}) PlaceholderAction {
	return func() utils.DocxWriter {
		return func(fileContent string) string {
			updatedContent := fileContent

			// Process each loop key in the data map
			for loopKey, loopData := range data {
				loopItems, ok := loopData.([]map[string]string)
				if !ok {
					continue // If not a slice of maps, skip this key
				}

				// Construct loop markers
				startMarker := "{{#each " + loopKey + "}}"
				endMarker := "{{/each}}"

				// Find the loop section in the updatedContent
				startLoop := strings.Index(updatedContent, startMarker)
				endLoop := strings.Index(updatedContent, endMarker) + len(endMarker)
				if startLoop == -1 || endLoop == -1 {
					continue // Loop markers not found, continue with next
				}

				// Extract the template section for the loop
				loopTemplate := updatedContent[startLoop+len(startMarker) : endLoop-len(endMarker)]

				// Generate content for each item in the loop
				var result strings.Builder
				for _, item := range loopItems {
					iterationContent := loopTemplate
					for key, value := range item {
						placeholder := "{{" + key + "}}"
						iterationContent = strings.ReplaceAll(iterationContent, placeholder, value)
					}
					result.WriteString(iterationContent)
				}

				// Replace the original loop section with generated content
				updatedContent = updatedContent[:startLoop] + result.String() + updatedContent[endLoop:]
			}

			return updatedContent
		}
	}
}
