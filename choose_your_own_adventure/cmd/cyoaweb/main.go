package main

import (
	"flag"
	"fmt"
	"os"

	cyoa "github.com/Joe-Bresee/gophercises/choose_your_own_adventure"
)

func main() {
	file := flag.String("file", "gopher.json", "JSON file containing cyoa story")
	flag.Parse()
	fmt.Printf("using the story in %s", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
