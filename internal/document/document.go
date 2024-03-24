package document

import (
	"archive/zip"
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliamerj/docxer/internal/template"
	"github.com/aliamerj/docxer/internal/utils"
)

//go:embed template/*
var documentXml embed.FS

func CreateNewDocx(dirPath string, title string, body string) (string, error) {
	outputFilePath := filepath.Join(dirPath, "new_file.docx")
	file, err := os.Create(outputFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	if err := template.CreateDocxTemplate(zipFile); err != nil {
		return "", err
	}
	docxer := docxWriter(title, body)
	if err := utils.CreateDocx(documentXml, zipFile, docxer); err != nil {
		return "", err
	}
	return outputFilePath, err
}
func docxWriter(title string, body string) utils.DocxWriter {

	return func(fileContent string) string {
		updatedTemplate := strings.Replace(fileContent, "{{TITLE}}", title, -1)
		updatedTemplate = strings.Replace(updatedTemplate, "{{BODY}}", body, -1)
		return updatedTemplate
	}
}
