package main

import (
	"fmt"
	"github.com/aliamerj/docxer/pkg/docxer"
	"log"
)

func main() {
	// with basic docx
	dx := docxer.NewDocx()
	dx.Title = "My Title"
	dx.Body = "This is the document body."
	path, err := dx.CreateNewDocx(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path)

	// with Markdown
	text := `
#$# New 
# Heading level 1 
## Heading level 2 
### Heading level 3 
#### Heading level 4 
##### Heading level 5
###### Heading level 6 
I really like using Markdown.
  `
	path2, err2 := docxer.CreateMarkdownDocx(".", text)
	if err != nil {
		log.Fatal(err2)
	}
	fmt.Println(path2)
}
