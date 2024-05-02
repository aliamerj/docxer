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

func Placeholder(filePath string) *holder {
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
	action := placeholder.TextPlaceholderWriter(replacements)
	err := placeholder.UpdateDocx(h.filePath, action)
	if err != nil {
		return err
	}
	return nil
}
func (h *holder) Loop(loop map[string]interface{}) error {
	dirPath := filepath.Dir(h.filePath)
	if err := utils.ValidateFilePath(dirPath); err != nil {
		return err
	}
	action := placeholder.LoopPlaceholderWriter(loop)
	err := placeholder.UpdateDocx(h.filePath, action)
	if err != nil {
		return err
	}

	return nil
}

func (h *holder) Image(imagePath, imageId string, widthEMU, heightEMU int) error {
	dirPath := filepath.Dir(h.filePath)
	if err := utils.ValidateFilePath(dirPath); err != nil {
		return err
	}
	action := placeholder.ImagePlaceholderWriter(imagePath, imageId, widthEMU, heightEMU)
	err := placeholder.UpdateDocx(h.filePath, action)
	if err != nil {
		return err
	}

	return nil

}
