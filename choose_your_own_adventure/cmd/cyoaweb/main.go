package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/Joe-Bresee/gophercises/choose_your_own_adventure"
)

func main() {
	port := flag.Int("port", 3000, "port to run the webapp story")
	file := flag.String("file", "gopher.json", "JSON file containing cyoa story")
	flag.Parse()
	fmt.Printf("using the story in %s\n", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting server on %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
