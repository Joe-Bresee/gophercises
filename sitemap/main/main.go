package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	linkparser "github.com/Joe-Bresee/gophercises/html_link_parser"
)

// types
const xmlns = "http://sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"Loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

// main
func main() {

	urlName := flag.String("url", "https://gophercises.com", "name of site you want to map")
	maxDepth := flag.Int("int", 10, "maxdepth of sitemap")
	flag.Parse()

	/*
	   1. GET the webpage
	   2. parse all the links on the page
	   3. build proper (full) urls with our links
	   4. filter out any links w/ a diff domain
	   5. Find all pages (BFS)
	   6. print out XML
	*/
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
	return hrefs(resp.Body, base)
}

// Finds hrefs in html, returns full link not just relative. must fall within base domain
func hrefs(body io.Reader, base string) []string {
	links, err := linkparser.Parse(body)
	if err != nil {
		log.Fatal(err)
	}
	var ret []string
	for _, link := range links {
		if strings.HasPrefix(link.Href, "/") {
			ret = append(ret, base+link.Href)
		} else if strings.HasPrefix(link.Href, "http://"+base) {
			ret = append(ret, link.Href)
		} else {
			continue
		}
	}
	return ret
}
