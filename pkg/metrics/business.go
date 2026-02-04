package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 게시글 수
	PostsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "board_posts_total",
		Help: "Total number of posts",
	})

	// 게시글 생성 카운터
	PostsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "board_posts_created_total",
		Help: "Total number of posts created",
	})

	// 댓글 수
	CommentsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "board_comments_total",
		Help: "Total number of comments",
	})

	// 사용자 로그인 카운터
	UserLogins = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "board_user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"status"}, // success, failure
	)

	// 데이터베이스 쿼리 시간
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "board_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"operation", "table"},
	)
)
