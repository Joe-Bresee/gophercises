package urlshort

import (
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Return anon function to handle the http request for when the request is performed.
	return func(w http.ResponseWriter, r *http.Request) {
		// no looping required since pathsToUrls is a map, simple lookup
		dest, found := pathsToUrls[r.URL.Path]
		if found {
			// perform redirection to destination url with appropriate http status code
			http.Redirect(w, r, dest, http.StatusFound)
		} else {
			// if not, return fallback default mux serving
			fallback.ServeHTTP(w, r)
		}
	}
}

// yamlhandler will parse the provided yaml and then return
// an http.handlerfunc (which also implements http.handler)
// that will attempt to map any paths to their corresponding
// url. if the path is not provided in the yaml, then the
// fallback http.handler will be called instead.
//
// yaml is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// the only errors that can be returned all related to having
// invalid yaml data.
//
// see maphandler to create a similar http.handlerfunc via
// // a mapping of paths to urls.
// func yamlhandler(yml []byte, fallback http.handler) (http.handlerfunc, error) {

// 	// struct for yaml parsing
// 	type urlyamlmapping struct {
// 		path string `yaml:"path"`
// 		URL  string `yaml:"url"`
// 	}

// 	// yaml map slice
// 	var u []urlYamlMapping

// 	// unmarshal the myl []byte into the urlYamlMapping
// 	err := yaml.Unmarshal([]byte(yml), &u)
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}

// 	// efficiently handle urlshort lookups by changing data struct into a map
// 	pathToUrl := make(map[string]string)
// 	for _, item := range u {
// 		pathToUrl[item.Path] = item.URL
// 	}

// 	func http.HandlerFunc(){
// 		// write httphandler
// 	}
// 	return nil, nil
// }
