package document

import (
	"archive/zip"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliamerj/docxer/internal/template"
)

//go:embed template/document.xml
var documentXml embed.FS

func CreateNewDocx(dirPath string, title string, body string) (string, error) {
	outputFilePath := filepath.Join(dirPath, "new_file.docx")
	// create new file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// create zip file
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	if err := template.CreateDocxTemplate(zipFile); err != nil {
		return "", err
	}
	fileXml, err := fs.ReadFile(documentXml, "template/document.xml")
	if err != nil {
		return "", err
	}
	newFile, err := zipFile.Create("word/document.xml")
	if err != nil {
		return "", err
	}
	updatedTemplate := strings.Replace(string(fileXml), "{{TITLE}}", title, -1)
	updatedTemplate = strings.Replace(updatedTemplate, "{{BODY}}", body, -1)

	if _, err := newFile.Write([]byte(updatedTemplate)); err != nil {
		return "", err
	}

	return outputFilePath, err

}
