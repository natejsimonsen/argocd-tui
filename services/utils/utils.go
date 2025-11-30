package utils

import (
	"regexp"
)

var tviewTagRegex = regexp.MustCompile(`\[.*?\]`)

func StripTags(text string) string {
	return tviewTagRegex.ReplaceAllString(text, "")
}

func GetTag(text string) string {
	return tviewTagRegex.FindString(text)
}
