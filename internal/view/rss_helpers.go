package view

import (
	"regexp"
	"strings"
)

var rssHTMLTagRe = regexp.MustCompile(`<[^>]*>`)

// RSSStripHTML strips HTML tags from a string.
func RSSStripHTML(s string) string {
	return strings.TrimSpace(rssHTMLTagRe.ReplaceAllString(s, ""))
}

// RSSTextTruncate truncates a string to maxLen runes, appending "..." if truncated.
func RSSTextTruncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
