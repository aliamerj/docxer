package main

import (
	"fmt"
	"log"

	"github.com/aliamerj/docxer/docxer"
)

func main() {
	dx := docxer.New()
	dx.Title = "My Title"
	dx.Body = "This is the document body."

	path, err := dx.CreateNewDocx(".")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path)
}
