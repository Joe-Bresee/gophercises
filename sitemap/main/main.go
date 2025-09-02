package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	linkparser "github.com/Joe-Bresee/gophercises/html_link_parser"
)

// types
const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

// main
func main() {

	urlFlag := flag.String("url", "https://gophercises.com", "name of site you want to map")
	maxDepth := flag.Int("depth", 10, "maxdepth of sitemap")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}

func bfs(urlStr string, maxDepth int) []string {
	// using struct for set maps uses less memory apparently than bool: Go doesn't have sets built-in
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {break}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue //if we have a link we've seen, skip we don't need to go over it again
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
				}
			}
		}
	ret := make([]string, 0, len(seen)) //pre-initialize ret slice for runtime speed improvement
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

// performs get req, hands body and url bas to href.
func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

// Finds hrefs in html, returns full link not just relative. must fall within base domain
func hrefs(body io.Reader, base string) []string {
	links, _ := linkparser.Parse(body)
	
	var ret []string
	for _, link := range links {
		if strings.HasPrefix(link.Href, "/") {
			ret = append(ret, base+link.Href)
		} else if strings.HasPrefix(link.Href, base) {
			ret = append(ret, link.Href)
		} else {
			continue
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
			}
		}
		return ret
	}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}