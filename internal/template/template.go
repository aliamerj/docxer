package template

import (
	"archive/zip"
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

func CreateDocxTemplate(zipFile *zip.Writer) error {
	return fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through templates: %w", err)
		}
		if d.IsDir() {
			return nil // Skip directories, as we only want to process files
		}

		// Prepare the ZIP file path based on whether it's a .rels file or not
		filePath := strings.TrimPrefix(path, "templates/")
		var zipPath string
		if strings.HasSuffix(filePath, ".rels") {
			zipPath = "_rels/" + filePath // Special handling for .rels files
		} else {
			zipPath = filePath
		}

		// Attempt to create a new file within the ZIP archive
		newFile, err := zipFile.Create(zipPath)
		if err != nil {
			return fmt.Errorf("error creating file '%s' in ZIP archive: %w", zipPath, err)
		}

		// Read and write the file content to the ZIP archive
		fileContent, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("error reading contents of '%s': %w", path, err)
		}

		if _, err := newFile.Write(fileContent); err != nil {
			return fmt.Errorf("error writing contents to '%s' in ZIP archive: %w", zipPath, err)
		}

		return nil
	})
}
