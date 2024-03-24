package docxer

import (
	"github.com/aliamerj/docxer/internal/document"
	"github.com/aliamerj/docxer/internal/markdown"
	"github.com/aliamerj/docxer/internal/utils"
)

type docxer struct {
	Title string
	Body  string
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
