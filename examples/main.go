package main

import (
	"fmt"
	"github.com/aliamerj/docxer/pkg/docxer"
	"log"
)

func main() {
	dx := docxer.NewDocx()
	dx.Title = "My Title"
	dx.Body = "This is the document body."
	path, err := dx.CreateNewDocx(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path)
}
