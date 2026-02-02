package router

import (
	"gorm-test/internal/handler"

	"github.com/gin-gonic/gin"
)

// Router 라우터
type Router struct {
	engine      *gin.Engine
	postHandler *handler.PostHandler
}

// NewRouter 생성자
func NewRouter(postHandler *handler.PostHandler) *Router {
	return &Router{
		engine:      gin.Default(),
		postHandler: postHandler,
	}
}

// Setup 라우트 설정
func (r *Router) Setup() *gin.Engine {
	// API 버전 그룹
	v1 := r.engine.Group("/api/v1")
	{
		// 게시글 라우트
		posts := v1.Group("/posts")
		{
			posts.GET("", r.postHandler.GetList)
			posts.POST("", r.postHandler.Create)
			posts.GET("/:id", r.postHandler.GetByID)
			posts.PUT("/:id", r.postHandler.Update)
			posts.DELETE("/:id", r.postHandler.Delete)
		}
	}

	// 헬스 체크
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
