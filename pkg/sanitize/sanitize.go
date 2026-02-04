package sanitize

import (
	"html"
	"regexp"
	"strings"
)

// HTML은 HTML 특수 문자를 이스케이프합니다
func HTML(s string) string {
	return html.EscapeString(s)
}

// StripTags는 모든 HTML 태그를 제거합니다
func StripTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// TrimSpace는 앞뒤 공백을 제거합니다
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// Normalize는 일반적인 정제를 수행합니다
func Normalize(s string) string {
	s = TrimSpace(s)
	s = StripTags(s)
	return s
}
