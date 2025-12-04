package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var tviewTagRegex = regexp.MustCompile(`\[.*?\]`)

func StripTags(text string) string {
	return tviewTagRegex.ReplaceAllString(text, "")
}

func GetTag(text string) string {
	return tviewTagRegex.FindString(text)
}

func HexToColor(hexStr string) tcell.Color {
	hex := strings.TrimPrefix(hexStr, "#")

	v, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		fmt.Printf("Error converting hex color %s: %v. Returning ColorDefault.\n", hexStr, err)
		return tcell.ColorDefault
	}

	return tcell.NewHexColor(int32(v))
}
