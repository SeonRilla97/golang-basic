# 미들웨어

요청과 응답 사이에 위치하는 코드 
 - 클라이언트의 요청이 핸들러에 도달하기 전에 먼저 처리, 핸들러가 응답을 보낸 후에 추가 작업 수행
 - 횡단 관심사(Cross-cutting Concerns)를 미들웨어로 분리


```yaml
# 타입
type HandlerFunc func(*Context)

func MyMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {

  // 요청 전 처리 (Pre-processing) : 핸들러가 실행되기 전에 수행할 작업

  c.Next() // 다음 핸들러 호출

  // 응답 후 처리 (Post-processing) : 핸들러가 실행된 후에 수행할 작업 
  }
}


미들웨어 메서드 & 설명
  c.Next()	다음 핸들러 실행
  c.Abort()	체인 중단 (현재 핸들러는 계속 실행)
  c.AbortWithStatus(code)	상태 코드와 함께 중단
  c.AbortWithStatusJSON(code, obj)	JSON 응답과 함께 중단

데이터 전달 및 받기
  c.Set()
  c.Get()
  c.MustGet("userID").(string) // 확실히 들어오지 않으면 패닉 발생함

미들웨어 사용


```
사용방법
```
## 전역
r := gin.New()
r.Use(gin.Logger(), gin.Recovery())

## 그룹
api := r.Group("/api")
api.Use(AuthMiddleware()) { 
    api.GET("/posts", GetPosts)
}

## 라우터
r.GET("/sensitive", AuthMiddleware(), SensitiveHandler)
```


웹 엔진 생성 방법
```
gin.Default()
gin.New() : Logger + Recovery 미들웨어 자동 추가

Logger MiddleWare
    gin.LoggerWithConfig() // API 로그 세밀한 조정
        LogFormatterParams 구조체 사용

gin.Recovery()
    gin.RecoveryWithWriter(w io.Writer) // 패닉 로그 출력 대상 변경  
    gin.CustomRecovery() // 복구 시 커스텀 동작 정의      
```


커스텀 미들웨어
```go
# gin.HandlerFunc 반환 함수로 작성
func MyMiddleware(config Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 전처리
        c.Next()
        // 후처리
    }
}
```

### 사용법 


```go

// 요청 별 고유 ID 부착
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 헤더에서 요청 ID 확인 (외부에서 전달된 경우)
        requestID := c.GetHeader(RequestIDHeader)

        // 없으면 새로 생성
        if requestID == "" {
            requestID = uuid.New().String()
        }

        // 컨텍스트와 응답 헤더에 설정
        c.Set("requestID", requestID)
        c.Header(RequestIDHeader, requestID)

        c.Next()
    }
}

// 요청 처리 시간을 측정하고 헤더로 반환
func Timing() gin.HandlerFunc {
    return func(c *gin.Context) {
    start := time.Now()
    
    c.Next()
    
    duration := time.Since(start)
    
    // 응답 헤더에 처리 시간 추가
    c.Header("X-Response-Time", fmt.Sprintf("%dms", duration.Milliseconds()))
    
    // 컨텍스트에도 저장 (로깅용)
    c.Set("latency", duration)
    }
}

// 메소드 오버라이드
// HTML- GET/POST만 지원 -> HTML 폼에서 PUT / DELETE 지원을 위한 미들웨어
func MethodOverride() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == http.MethodPost {
            // _method 파라미터 확인
            method := c.PostForm("_method")
        if method == "" {
            // X-HTTP-Method-Override 헤더 확인
            method = c.GetHeader("X-HTTP-Method-Override")
        }
        
        // 유효한 메서드로 오버라이드
        switch method {
        case http.MethodPut, http.MethodPatch, http.MethodDelete:
            c.Request.Method = method
        }
    }
    
    c.Next()
    }
}


// Conditional은 조건에 따라 미들웨어를 적용합니다
// 특정 경로는 비인증 처리 등
type SkipFunc func(c *gin.Context) bool
func Conditional(skipFunc SkipFunc, middleware gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        if skipFunc(c) {
            c.Next()
            return
    }
        middleware(c)
    }
}
```

프로젝트의 규모가 커지면 미들웨어도 모듈화하고 설정 가능하도록 만든다

```go
type Config struct {
    // 요청 ID 설정
    RequestID struct {
        Enabled bool
        Header  string
    }

    // 타이밍 설정
    Timing struct {
        Enabled bool
        Header  string
    }

    // 로깅 설정
    Logging struct {
        Enabled   bool
        SkipPaths []string
    }
}

func DefaultConfig() Config {
    return Config{
        RequestID: struct {
            Enabled bool
            Header  string
        }{
            Enabled: true,
            Header:  "X-Request-ID",
        },
        Timing: struct {
            Enabled bool
            Header  string
        }{
            Enabled: true,
            Header:  "X-Response-Time",
        },
        Logging: struct {
            Enabled   bool
            SkipPaths []string
        }{
            Enabled:   true,
            SkipPaths: []string{"/health", "/metrics"},
        },
    }
}

func Setup(r *gin.Engine, cfg Config) {
    if cfg.RequestID.Enabled {
        r.Use(RequestIDWithHeader(cfg.RequestID.Header))
    }
    
    if cfg.Timing.Enabled {
        r.Use(TimingWithHeader(cfg.Timing.Header))
    }
    
    if cfg.Logging.Enabled {
        r.Use(LoggingWithSkip(cfg.Logging.SkipPaths))
    }
}
```

미들웨어의 순서
```go
// 일반적인 순서
r.Use(
    // 1. 기본 설정 (항상 먼저)
    middleware.RequestID(),
    middleware.Timing(),

    // 2. 로깅 (요청 ID, 타이밍 정보 필요)
    middleware.Logger(),

    // 3. 보안 (로깅 후, 비즈니스 로직 전)
    middleware.Recovery(),
    middleware.SecureHeaders(),
    middleware.CORS(),

    // 4. 요청 제한 (인증 전)
    middleware.RateLimiter(),

    // 5. 인증/인가 (비즈니스 로직 직전)
    middleware.Auth(),
)

// 미들웨어 조합을 함수로 만들어 재사용
// middleware/chains.go
package middleware

import "github.com/gin-gonic/gin"

// CommonMiddlewares returns common middleware chain
func CommonMiddlewares() []gin.HandlerFunc {
return []gin.HandlerFunc{
RequestID(),
Timing(),
Logger(),
Recovery(),
}
}

// APIMiddlewares returns API-specific middleware chain
func APIMiddlewares() []gin.HandlerFunc {
return []gin.HandlerFunc{
CORS(),
RateLimiter(100, time.Minute),
SecureHeaders(),
}
}

// AuthMiddlewares returns authentication middleware chain
func AuthMiddlewares() []gin.HandlerFunc {
return []gin.HandlerFunc{
Auth(),
}
}

r := gin.New()
r.Use(middleware.CommonMiddlewares()...)

api := r.Group("/api")
api.Use(middleware.APIMiddlewares()...)
{
// 공개 라우트
api.GET("/posts", postHandler.List)

// 인증 필요 라우트
authGroup := api.Group("")
authGroup.Use(middleware.AuthMiddlewares()...)
{
authGroup.POST("/posts", postHandler.Create)
}
}


// 동적 미들웨어
type DynamicMiddleware struct {
middlewares []gin.HandlerFunc
}

func NewDynamic() *DynamicMiddleware {
return &DynamicMiddleware{
middlewares: make([]gin.HandlerFunc, 0),
}
}

func (d *DynamicMiddleware) Add(m gin.HandlerFunc) {
d.middlewares = append(d.middlewares, m)
}

func (d *DynamicMiddleware) Handler() gin.HandlerFunc {
return func(c *gin.Context) {
// 등록된 미들웨어들을 순차 실행
for _, m := range d.middlewares {
m(c)
if c.IsAborted() {
return
}
}
c.Next()
}
}

```


