package middleware

// Redis를 이용한 분산 Rate Limit

//type RedisRateLimiter struct {
//	client *redis.Client
//	limit  int
//	window time.Duration
//}
//
//func NewRedisRateLimiter(client *redis.Client, limit int, window time.Duration) *RedisRateLimiter {
//	return &RedisRateLimiter{
//		client: client,
//		limit:  limit,
//		window: window,
//	}
//}
//
//func (rl *RedisRateLimiter) Middleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx := context.Background()
//		ip := c.ClientIP()
//		key := fmt.Sprintf("rate_limit:%s", ip)
//
//		// 현재 카운트 조회 및 증가
//		count, err := rl.client.Incr(ctx, key).Result()
//		if err != nil {
//			c.Next() // Redis 에러 시 허용
//			return
//		}
//
//		// 첫 요청이면 만료 시간 설정
//		if count == 1 {
//			rl.client.Expire(ctx, key, rl.window)
//		}
//
//		// 제한 초과 확인
//		if count > int64(rl.limit) {
//			ttl, _ := rl.client.TTL(ctx, key).Result()
//			c.Header("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())+1))
//			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
//				"success": false,
//				"error": gin.H{
//					"code":    "RATE_LIMITED",
//					"message": "요청이 너무 많습니다.",
//				},
//			})
//			return
//		}
//
//		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
//		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.limit-int(count)))
//
//		c.Next()
//	}
//}
