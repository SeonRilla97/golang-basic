package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 요청 수 카운터
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// 요청 지연 시간 히스토그램
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 현재 처리 중인 요청 수
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath() // 파라미터가 포함된 경로 (/posts/:id)
		if path == "" {
			path = "unknown"
		}

		// 처리 중인 요청 수 증가
		httpRequestsInFlight.Inc()

		start := time.Now()

		c.Next()

		// 처리 중인 요청 수 감소
		httpRequestsInFlight.Dec()

		// 메트릭 기록
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}
