package main

import (
	"fmt"
	"gorm-test/internal/auth"
	"gorm-test/internal/config"
	"gorm-test/internal/database"
	"gorm-test/internal/handler"
	"gorm-test/internal/repository"
	"gorm-test/internal/router"
	"gorm-test/internal/service"
	"gorm-test/middleware"
	"gorm-test/pkg/notify"
	"gorm-test/pkg/sentry"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 설정 로드
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Sentry 초기화
	if err := sentry.Init(cfg.Sentry.Dsn); err != nil {
		log.Fatal(err)
	}

	// Gin 모드 설정
	gin.SetMode(cfg.Server.Mode)

	// 데이터베이스 연결
	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	// 의존성 주입
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo, cfg)
	postHandler := handler.NewPostHandler(postService)

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo, postRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	userRepo := repository.NewUserRepository(db)
	tokenService := auth.NewTokenService("secreykkkkkkkkkkkkey", 1, 2)
	passwordService := auth.NewPasswordService()
	authService := service.NewAuthService(userRepo, passwordService, tokenService)
	authHandler := handler.NewAuthHandler(authService) // 라우터 설정
	r := router.NewRouter(postHandler, commentHandler, authHandler)

	corsConfig := middleware.CORSConfig{
		Debug: cfg.Server.Env == "development",
		AllowedOrigins: []string{
			"https://example.com",
			"https://admin.example.com",
		},
	}

	// Recovery 설정
	slackWebhook := "test"
	recoveryConfig := middleware.RecoveryConfig{
		EnableStackTrace: true,
		NotifyFunc: func(err any, stack string) {
			notify.SendToSentry(err, stack)
			notify.SendToSlack(slackWebhook)(err, stack)
		},
	}
	r.Use(middleware.CORS(corsConfig))
	r.Use(middleware.RateLimiter(10, 20))
	r.Use(middleware.SecureHeaders(middleware.DefaultSecureConfig(cfg)))

	ipLimiter := middleware.NewIPRateLimiter(5, 10, time.Hour)
	r.Use(ipLimiter.Middleware())

	if cfg.Server.Env == "development" {
		r.Use(
			middleware.BodyLogging(1024),
			middleware.RequestID(),
			middleware.Prometheus(), // 메트릭 수집
			middleware.Logging(),
			middleware.Recovery(recoveryConfig),
			middleware.ErrorHandler(), // 에러 핸들러
		) // 1KB 제한
	}

	engine := r.Setup()

	// 서버 시작
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("서버 시작: http://localhost%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}
