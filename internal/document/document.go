package document

import (
	"archive/zip"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliamerj/docxer/internal/template"
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
	if err := createDocx(zipFile, title, body); err != nil {
		return "", err
	}
	return outputFilePath, err
}

func createDocx(zipFile *zip.Writer, title, body string) error {
	return fs.WalkDir(documentXml, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through template: %w", err)
		}
		if d.IsDir() {
			return nil
		}
		filePath := strings.TrimPrefix(path, "template/")
		newFile, err := zipFile.Create("word/" + filePath)
		if err != nil {
			return fmt.Errorf("error creating file '%s' in ZIP archive: %w", filePath, err)
		}
		fileContent, err := fs.ReadFile(documentXml, path)
		if err != nil {
			return fmt.Errorf("error reading contents of '%s': %w", path, err)
		}
		updatedTemplate := strings.Replace(string(fileContent), "{{TITLE}}", title, -1)
		updatedTemplate = strings.Replace(updatedTemplate, "{{BODY}}", body, -1)

		if _, err := newFile.Write([]byte(updatedTemplate)); err != nil {
			return err
		}
		return nil
	})
}
