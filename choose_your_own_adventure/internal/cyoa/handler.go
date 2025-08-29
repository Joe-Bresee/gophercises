package cyoa

import "net/http"

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "" || path == "/" {
		path = "/intro"
	}

	chapterKey := path[1:]

	chapter, ok := h.s[chapterKey]
	if !ok {
		http.Error(w, "Chapter not found", http.StatusNotFound)
		return
	}

	err := tpl.Execute(w, chapter)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}
