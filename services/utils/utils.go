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

func HexToColor(hexStr string, def tcell.Color) tcell.Color {
	if hexStr == "" {
		return def
	}

	hex := strings.TrimPrefix(hexStr, "#")

	v, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		fmt.Printf("Error converting hex color %s: %v. Returning ColorDefault.\n", hexStr, err)
		return tcell.ColorDefault
	}

	userColor := tcell.NewHexColor(int32(v))
	return userColor
}

func GetContrastColor(c tcell.Color) tcell.Color {
	r, g, b := c.RGB()

	luminance := (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) / 255

	if luminance > 0.5 {
		return tcell.ColorBlack
	}
	return tcell.ColorWhite
}
