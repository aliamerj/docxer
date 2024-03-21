package document

import (
	"archive/zip"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

func CreateNewDocx(dirPath string, title string, Body string) (string, error) {
	if err := validateFilePath(dirPath); err != nil {
		return "", err
	}

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

	if err := createDocxTemplete(zipFile, title, Body); err != nil {
		return "", err
	}

	return outputFilePath, err

}

func createDocxTemplete(zipFile *zip.Writer, title string, body string) error {

	if err := fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		filePath := strings.Replace(path, "templates/", "", -1)

		if filePath == ".rels" {
			fmt.Println(path)
			newFile, err := zipFile.Create("_rels/" + filePath)
			if err != nil {
				return err
			}
			fileContent, err := fs.ReadFile(templateFS, path)
			if err != nil {
				return err
			}

			if _, err := newFile.Write(fileContent); err != nil {
				return err
			}
		}

		newFile, err := zipFile.Create(filePath)
		if err != nil {
			return err
		}
		fileContent, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return err
		}
		if filePath == "word/document.xml" {
			fmt.Println(filePath)
			updatedTemplate := strings.Replace(string(fileContent), "{{TITLE}}", title, -1)
			updatedTemplate = strings.Replace(updatedTemplate, "{{BODY}}", body, -1)

			if _, err := newFile.Write([]byte(updatedTemplate)); err != nil {
				return err
			}
			return nil

		}

		if _, err := newFile.Write(fileContent); err != nil {
			return err
		}
		return nil

	}); err != nil {
		return err
	}
	return nil
}

func validateFilePath(dirPath string) error {
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
