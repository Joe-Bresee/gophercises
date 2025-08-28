package cyoa

import "net/http"

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
