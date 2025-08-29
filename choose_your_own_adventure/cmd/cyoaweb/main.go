package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Joe-Bresee/gophercises/choose_your_own_adventure/internal/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "port to run the webapp story")
	file := flag.String("file", "./data/gopher.json", "JSON file containing cyoa story")
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
