package template

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateDocxTemplate(t *testing.T) {
	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	// Assuming your CreateDocxTemplate function and context setup is corrected and available
	err := CreateDocxTemplate(zipWriter)
	if err != nil {
		t.Fatalf("CreateDocxTemplate failed: %v", err)
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("Failed to close zip writer: %v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	if err != nil {
		t.Fatalf("Failed to read zip content: %v", err)
	}

	foundRels := false
	foundNonRels := false

	for _, file := range zipReader.File {
		if strings.Contains(file.Name, ".rels") {
			if !strings.HasPrefix(file.Name, "_rels/") {
				t.Errorf("RELs file %s is not in _rels directory", file.Name)
			} else {
				foundRels = true
			}
		} else {
			foundNonRels = true
		}
	}

	if !foundRels {
		t.Errorf("No .rels files found in _rels directory")
	}

	if !foundNonRels {
		t.Errorf("No non-.rels files found outside _rels directory")
	}
}
