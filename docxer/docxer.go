package docxer

import "github.com/aliamerj/docxer/pkg/document"

type docxer struct {
	Title string
	Body  string
}

func NewDocx() *docxer {
	return &docxer{}
}
func (d *docxer) CreateNewDocx(filePath string) (string, error) {
 if err := validateFilePath(filePath); err != nil {
		return "", err
	}


	path, err := document.CreateNewDocx(filePath, d.Title, d.Body)
	if err != nil {
		return "", err
	}
	return path, nil
}
func NewMarkdownDocx(filePath string, markdown string) (string, error){
 if err := validateFilePath(filePath); err != nil {
		return "", err
	}
  // todo
  return "", nil
}


