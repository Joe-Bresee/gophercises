package cyoa

import (
	"encoding/json"
	"io"
)

func JsonStory(r io.Reader) (Story, error) {
	// pass in io.reader to newdecoder instead of usign bytes unmarshalling
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

// generated from gopher.json via https://mholt.github.io/json-to-go/
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
