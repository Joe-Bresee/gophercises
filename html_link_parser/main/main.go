package main

import (
	"fmt"
	"io"
	"log"
	"os"

	linkparser "github.com/Joe-Bresee/gophercises/html_link_parser"
)

func main() {
	file, err := os.Open("./html/example_href.html")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	r := io.Reader(file)
	links, err := linkparser.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range links {
		fmt.Printf("href: %s, Text: %q\n", link.Href, link.Text)
	}
}
