package main

import "regexp"

var nonDigits = regexp.MustCompile(`\D+`)

func normalize(s string) string {
	return nonDigits.ReplaceAllString(s, "")
}
