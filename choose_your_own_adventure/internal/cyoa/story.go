package cyoa

import (
	"encoding/json"
	"io"
	"text/template"
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

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("internal/templates/story.html"))
}
