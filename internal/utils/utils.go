package utils

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func ValidateFilePath(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	// Check if the path exists and is a directory
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the specified directory does not exist: %s", dirPath)
		}
		return fmt.Errorf("error checking the specified directory: %v", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("the specified path is not a directory: %s", dirPath)
	}

	return nil
}

type DocxWriter func(string) string

func CreateDocx(documentXml fs.FS, zipFile *zip.Writer, docxWriter DocxWriter) error {
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
		updatedTemplate := docxWriter(string(fileContent))
		if _, err := newFile.Write([]byte(updatedTemplate)); err != nil {
			return err
		}
		return nil
	})
}
