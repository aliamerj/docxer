package main

import (
	"fmt"
	"log"

	"github.com/aliamerj/docxer/pkg/docxer"
)

func main() {
	//	with basic docx
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
	I really like using Markdown in Docxer.
	you can make the text **Bold** , *Italic* or ***both***
	  `
	path2, err2 := docxer.CreateMarkdownDocx(".", text)
	if err != nil {
		log.Fatal(err2)
	}
	fmt.Println(path2)

	dox := docxer.Placeholder("./new_file.docx")
	dox.Text(map[string]string{"name": "Ali", "job": "Software Engineer"})
	replacements := map[string]string{
		"INVOICE_NUMBER": "123456",
		"DATE":           "2024-04-30",
		"TOTAL":          "150.00",
	}
	if err := dox.Text(replacements); err != nil {
		log.Fatal(err)
	}

	items := map[string]interface{}{
		"items": []map[string]string{
			{"NAME": "Product 1", "QUANTITY": "2", "PRICE": "30.00"},
			{"NAME": "Product 2", "QUANTITY": "1", "PRICE": "90.00"},
			{"NAME": "Product 3", "QUANTITY": "1", "PRICE": "90.00"},
			{"NAME": "Product 4", "QUANTITY": "1", "PRICE": "90.00"},
		},
		"products": []map[string]string{
			{"NAME": "Product x1", "QUANTITY": "200", "PRICE": "30.00"},
			{"NAME": "Product x2", "QUANTITY": "100", "PRICE": "90.00"},
			{"NAME": "Product x2x3", "QUANTITY": "100", "PRICE": "90.00"},
		},
	}
	dox.Loop(items)

}
