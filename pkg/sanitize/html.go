package sanitize

import "github.com/microcosm-cc/bluemonday"

var (
	// 모든 HTML 태그 제거
	strictPolicy = bluemonday.StrictPolicy()

	// 안전한 태그만 허용 (링크, 포맷팅)
	ugcPolicy = bluemonday.UGCPolicy()
)

// Strict는 모든 HTML을 제거합니다
func Strict(s string) string {
	return strictPolicy.Sanitize(s)
}

// UGC는 사용자 생성 콘텐츠용으로 안전한 HTML만 허용합니다
func UGC(s string) string {
	return ugcPolicy.Sanitize(s)
}

// 커스텀 정책 생성
func NewCustomPolicy() *bluemonday.Policy {
	p := bluemonday.NewPolicy()

	// 허용할 태그와 속성 지정
	p.AllowElements("p", "br", "b", "i", "u", "strong", "em")
	p.AllowElements("ul", "ol", "li")
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")

	// 링크 허용 (rel 속성 강제)
	p.AllowAttrs("href").OnElements("a")
	p.RequireNoReferrerOnLinks(true)

	// 이미지 허용 (특정 도메인만)
	p.AllowAttrs("src", "alt").OnElements("img")
	p.AllowURLSchemes("https")

	return p
}
