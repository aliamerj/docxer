package docxer

import (
	"github.com/aliamerj/docxer/internal/document"
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
 if err :=utils.ValidateFilePath(filePath); err != nil {
		return "", err
	}


	path, err := document.CreateNewDocx(filePath, d.Title, d.Body)
	if err != nil {
		return "", err
	}
	return path, nil
}
