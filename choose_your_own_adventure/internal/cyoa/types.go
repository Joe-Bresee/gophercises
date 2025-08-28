package cyoa

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
