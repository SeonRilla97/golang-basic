package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowedOrigins []string
	Debug          bool
}

func CORS(cfg CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // 허용할 HTTP 메서드
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", // 허용할 헤더
		},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID", "X-Response-Time"}, // 응답에서 노출할 헤더
		AllowCredentials: true,                                                          // 쿠키 허용 여부
		MaxAge:           12 * time.Hour,                                                // Preflight 결과 캐시 시간
	}

	if cfg.Debug {
		// 개발 환경: 모든 Origin 허용
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = false // AllowAllOrigins와 함께 사용 불가
	} else {
		// 프로덕션: 지정된 Origin만 허용
		corsConfig.AllowOrigins = cfg.AllowedOrigins // 허용할 도메인 목록입니다. 와일드카드를 사용할 수 있지만, AllowCredentials가 true면 *를 사용할 수 없습니다.
	}

	return cors.New(corsConfig)
}
