package docxer

import (
	"path/filepath"

	"github.com/aliamerj/docxer/internal/document"
	"github.com/aliamerj/docxer/internal/markdown"
	"github.com/aliamerj/docxer/internal/placeholder"
	"github.com/aliamerj/docxer/internal/utils"
)

type docxer struct {
	Title string
	Body  string
}
type holder struct {
	filePath string
}

func Replace(filePath string) *holder {
	return &holder{filePath: filePath}
}

func NewDocx() *docxer {
	return &docxer{}
}
func (d *docxer) CreateNewDocx(filePath string) (string, error) {
	if err := utils.ValidateFilePath(filePath); err != nil {
		return "", err
	}
	path, err := document.CreateNewDocx(filePath, d.Title, d.Body)
	if err != nil {
		return "", err
	}
	return path, nil
}

func CreateMarkdownDocx(filePath string, markdownText string) (string, error) {
	if err := utils.ValidateFilePath(filePath); err != nil {
		return "", err
	}

	path, err := markdown.CreateMarkdownDocx(filePath, markdownText)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (h *holder) Text(replacements map[string]string) error {
	dirPath := filepath.Dir(h.filePath)
	if err := utils.ValidateFilePath(dirPath); err != nil {
		return err
	}
	err := placeholder.UpdateDocx(h.filePath, replacements)
	if err != nil {
		return err
	}
	return nil
}
