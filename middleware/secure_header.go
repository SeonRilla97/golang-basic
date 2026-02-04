package middleware

import (
	"gorm-test/internal/config"

	"github.com/gin-gonic/gin"
)

/*
CSP 설정

// API 서버 (리소스 로딩 없음)
cfg.ContentSecurityPolicy = "default-src 'none'; frame-ancestors 'none'"

// 웹 앱 서버
cfg.ContentSecurityPolicy = strings.Join([]string{
"default-src 'self'",
"script-src 'self' 'unsafe-inline' https://cdn.example.com",
"style-src 'self' 'unsafe-inline'",
"img-src 'self' data: https:",
"font-src 'self' https://fonts.gstatic.com",
"connect-src 'self' https://api.example.com",
"frame-ancestors 'none'",
}, "; ")

 // gin-helmet은 빠르게 기본 보안 헤더를 붙이고 싶을 때 편리한 라이브러리 -> 현재 Secure 헤더를 대체함


배포 후 보안 헤더를 점검할 수 있는 도구

Security Headers	https://securityheaders.com
Mozilla Observatory	https://observatory.mozilla.org
SSL Labs	https://www.ssllabs.com/ssltest

*/

type SecureConfig struct {
	// HSTS 설정
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool

	// 프레임 옵션
	FrameOption string // DENY, SAMEORIGIN

	// CSP 설정
	ContentSecurityPolicy string

	// Referrer 정책
	ReferrerPolicy string

	// 개발 모드 (일부 헤더 비활성화)
	Debug bool
}

func DefaultSecureConfig(config *config.Config) SecureConfig {
	return SecureConfig{
		HSTSMaxAge:            31536000, // 1년
		HSTSIncludeSubdomains: true,
		FrameOption:           "DENY",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		Debug:                 config.Server.Env == "development",
	}
}

func SecureHeaders(cfg SecureConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// MIME 스니핑 방지
		c.Header("X-Content-Type-Options", "nosniff")

		// XSS 필터 활성화
		c.Header("X-XSS-Protection", "1; mode=block")

		// 클릭재킹 방지
		c.Header("X-Frame-Options", cfg.FrameOption)

		// Referrer 정책
		c.Header("Referrer-Policy", cfg.ReferrerPolicy)

		// 권한 정책
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (프로덕션에서만)
		if !cfg.Debug {
			hsts := "max-age=" + string(rune(cfg.HSTSMaxAge))
			if cfg.HSTSIncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hsts)
		}

		// CSP (설정된 경우)
		if cfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}

		// 캐시 제어 (API 응답)
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")

		c.Next()
	}
}
