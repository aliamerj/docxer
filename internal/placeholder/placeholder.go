package placeholder

import (
	"archive/zip"
	"io"
	"os"
	"strings"

	"github.com/aliamerj/docxer/internal/utils"
)

type Replace struct {
	content string
	links   string
	images  map[string]string
}

func UpdateDocx(filePath string, replacements map[string]string) error {
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
	docxer := docxPlaceholderWriter(replacements)

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
func docxPlaceholderWriter(replacements map[string]string) utils.DocxWriter {
	return func(fileContent string) string {
		updatedTemplate := fileContent
		for placeholder, replacement := range replacements {
			updatedTemplate = strings.ReplaceAll(updatedTemplate, "{{"+placeholder+"}}", replacement)
		}
		return updatedTemplate
	}
}
