package router

import (
	"gorm-test/internal/auth"
	"gorm-test/internal/handler"
	"gorm-test/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router 라우터
type Router struct {
	engine         *gin.Engine
	postHandler    *handler.PostHandler
	commentHandler *handler.CommentHandler
	authHandler    *handler.AuthHandler
}

// NewRouter 생성자
func NewRouter(postHandler *handler.PostHandler, commentHandler *handler.CommentHandler, authHandler *handler.AuthHandler,
) *Router {
	return &Router{
		engine:         gin.Default(),
		postHandler:    postHandler,
		commentHandler: commentHandler,
	}
}

// Use 전역 미들웨어 등록
func (r *Router) Use(middleware ...gin.HandlerFunc) {
	r.engine.Use(middleware...)
}

// Setup 라우트 설정
func (r *Router) Setup(tokenService *auth.TokenService) *gin.Engine {
	// API 버전 그룹

	v1 := r.engine.Group("/api/v1")

	{

		//관리자 라우트
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(tokenService))
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", r.authHandler.Signup)
			admin.DELETE("/users/:id", r.authHandler.Signup)
			admin.GET("/stats", r.authHandler.Signup)
		}
		{

		// 인증 라우트
			authGroup := v1.Group("/api/auths")
		{
			authGroup.POST("/signup", r.authHandler.Signup)
			authGroup.POST("/login", r.authHandler.Login)
			authGroup.POST("/refresh", r.authHandler.RefreshToken)
			authGroup.POST("/logout", middleware.AuthMiddleware(tokenService, tokenStore), authHandler.Logout)
		}

		// 게시글 라우트 (비인증)
		postsPublic := v1.Group("/posts")
		postsPublic.Use(middleware.RateLimiter(5.0/60.0, 5)) // posts 는 분당 5회 Rate Limit 적용
		{
			// 댓글 라우트
			postsPublic.GET("/:postId/comments", r.commentHandler.GetByPostID)
		}
		// 게시글 라우트 (선택적 인증)
		postsOptional := v1.Group("/posts")
		postsOptional.Use(middleware.OptionalAuthMiddleware(tokenService))
		{
			postsOptional.GET("", r.postHandler.GetList)
			postsOptional.GET("/:id", r.postHandler.GetByID)
		}
		// 게시글 라우트 (인증)
		postsProtected := v1.Group("/posts")
		postsProtected.Use(middleware.AuthMiddleware(tokenService))

		{
			postsProtected.POST("", r.postHandler.Create)
			postsProtected.PUT("/:postId", r.postHandler.Update)
			postsProtected.DELETE("/:postId", r.postHandler.Delete)
			postsProtected.GET("/cursor", r.postHandler.GetListByCursor)
			// 댓글 라우트
			postsProtected.POST("/:postId/comments", r.commentHandler.Create)
			postsProtected.PUT("/:postId/comments/:commentId", r.commentHandler.Update)
			postsProtected.DELETE("/:postId/comments/:commentId", r.commentHandler.Delete)
		}
	}

	// 메트릭 엔드포인트 (인증 없이 접근 가능)
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 헬스 체크
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
