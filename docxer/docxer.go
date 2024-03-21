package docxer

import "github.com/aliamerj/docxer/pkg/document"

type docxer struct {
	Title string
	Body  string
}

func New() *docxer {
	return &docxer{}
}
func (d *docxer) CreateNewDocx(filePath string) (string, error) {
	path, err := document.CreateNewDocx(".", d.Title, d.Body)
	if err != nil {
		return "", err
	}
	return path, nil
}
