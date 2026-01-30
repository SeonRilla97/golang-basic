package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default() // gin.New() 미들웨어 없는 빈 라우터 - 학습 시 Default 사용

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.Run(":8080")
}
