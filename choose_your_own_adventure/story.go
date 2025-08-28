package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

// handle http w tpl which 'Must' template
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}

type Story map[string]Chapter

// generated from gopher.json via https://mholt.github.io/json-to-go/
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

var tpl *template.Template

func init() {
	tpl := template.Must(template.New("").Parse(defaultHandlerTmpl))
	fmt.Printf("%d", tpl)
}

// simple html tmplt
var defaultHandlerTmpl = `
<!DOCTYPE html>
    <html>
        <head>
            <meta charset="utf-8">
            <title>Choose Your Own Adventure</title>
        </head>
        <body>
            <!-- using html/template for dynamic html rendering -->
            <h1>{{ .Title }}</h1>
            {{range .Paragraphs}}
            <p>{{.}}</p>>
            {{end}}
            <ul>
                {{range .Options}}
                <li><a href="/{{ .Chapter }}">{{.Text}}</a></li>
                {{end}}
            </ul>
        </body>
    </html>`
