package middleware

import (
	"gorm-test/pkg/apperror"
	"gorm-test/pkg/logger"
	"gorm-test/pkg/problem"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 핸들러 실행

		// 에러가 없으면 종료
		if len(c.Errors) == 0 {
			return
		}

		// 마지막 에러 처리
		err := c.Errors.Last().Err
		instance := c.Request.URL.Path

		log := logger.FromGin(c)

		var pd *problem.Detail
		if appErr, ok := apperror.AsAppError(err); ok {
			pd = problem.FromAppError(appErr, instance)

			// 로깅
			attrs := []any{
				"error_code", appErr.Code,
				"http_status", appErr.HTTPStatus,
				"detail", appErr.Message,
			}
			if appErr.Err != nil {
				attrs = append(attrs, "error", appErr.Err)
			}
			log.Error(appErr.Message, attrs...)
		} else {
			pd = problem.InternalError(instance)
			log.Error("unhandled error", "error", err)
		}

		// 프로덕션 환경에서는 Detail 숨김
		if !gin.IsDebugging() {
			pd.Detail = ""
		}

		// RFC 7807 Content-Type 설정
		c.Header("Content-Type", problem.ContentType)
		c.JSON(pd.Status, pd)
	}
}
