package placeholder

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliamerj/docxer/internal/utils"
)

type PlaceholderAction func(zipWriter *zip.Writer, zipReader *zip.ReadCloser) utils.DocxWriter

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
	docxer := action(zipWriter, zipReader)

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
	return func(zipWriter *zip.Writer, zipReader *zip.ReadCloser) utils.DocxWriter {
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
	return func(zipWriter *zip.Writer, zipReader *zip.ReadCloser) utils.DocxWriter {
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

func ImagePlaceholderWriter(imagePath, placeholderID string, widthEMU, heightEMU int) PlaceholderAction {
	return func(zipWriter *zip.Writer, zipReader *zip.ReadCloser) utils.DocxWriter {
		// Prepare the image and add it to the ZIP
		imgData, err := os.ReadFile(imagePath)
		if err != nil {
			fmt.Println("Error reading image file:", err)
			return nil
		}
		imgFilename := fmt.Sprintf("image_%s%s", placeholderID, filepath.Ext(imagePath))
		imgPath := "word/media/" + imgFilename

		//1 Add image file to the DOCX package
		imageWriter, err := zipWriter.Create(imgPath)
		if err != nil {
			fmt.Println("Error creating zip entry for image:", err)
			return nil
		}
		_, err = imageWriter.Write(imgData)
		if err != nil {
			fmt.Println("Error writing image data to zip:", err)
			return nil
		}

		// Relationship ID should be uniquely generated; this is a simple placeholder.
		relID := "rId" + placeholderID

		//2 Add Relationship in document.xml.rels
		//TODO

		return func(fileContent string) string {
			imageXML := fmt.Sprintf(`<w:p><w:r><w:drawing><wp:inline width="%d" height="%d"><wp:extent cx="%d" cy="%d"/><wp:effectExtent l="0" t="0" r="0" b="0"/><wp:docPr id="%s" name="Picture 1"/><wp:cNvGraphicFramePr><a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/></wp:cNvGraphicFramePr><a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture"><pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture"><pic:nvPicPr><pic:cNvPr id="0" name=""/><pic:cNvPicPr/></pic:nvPicPr><pic:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill><pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic></a:graphicData></a:graphic></wp:inline></w:drawing></w:r></w:p>`, widthEMU, heightEMU, widthEMU, heightEMU, placeholderID, relID, widthEMU, heightEMU)
			return strings.ReplaceAll(fileContent, "{{%img%"+placeholderID+"}}", imageXML)
		}
	}
}
