Gin
 - 라우팅, 미들웨어, 요청 파싱, 에러 처리 쉽게 처리 가능


go get -u github.com/gin-gonic/gin



메서드	용도	예시
```
GET	리소스 조회	사용자 정보 조회
POST	리소스 생성	새 사용자 등록
PUT	리소스 전체 수정	사용자 정보 전체 업데이트
PATCH	리소스 부분 수정	사용자 이름만 변경
DELETE	리소스 삭제	사용자 삭제
```

- PATCH의 경우 포인터를 이용하여 빈값인지 아닌지 구분 가능
- Any : 모든 HTTP 메서드를 처리
- NoRoute : 정의되지 않은 요청 호출 시 처리
- Static() : 정적 파일 서빙
- StaticFS() : 임베디드 파일 시셑,ㅁ 커스텀 파일 시스템 사용 시
- StaticFile() : 단일 파일을 연결
- StaticFileFS() : 단일 파일 서빙
### URL 파라미터

- Path Parameter 
리소스 식별, 계층 구조 표현, 필수 값
  - [:변수명] -> c.Param("변수명")    : 변수 하나  
  - [*filepath] -> c.Param("filepath")  : 슬레시 포함 전체 경로 캡쳐

- Query String :
필터링, 정렬, 페이지네이션, 검색, 선택값
    - c.DefaultQuery("q","defaultValue") : ?q=golang
    - c.QueryArray("test") :  ?test=a&test=b&test=c
    - c.Query("q") : ?q=golang
    - c.QueryMap("config") : ?config[theme]=dark&config[lang]=ko

### 라우터 그룹 

```go
r := gin.Default() 
r.group("/api/v1"){
	r.GET("/users", func(c *gin.Context){ ... }) // /api/v1/users
	r.POST("/post", func(c *gin.Context){ ... }) // /api/v1/post
}
```
- API 버전 관리
- 중첩 그룹
- Group에 미들웨어 적용 시 하위에 모두 적용


### 바인딩

```
메서드	            데이터 소스	태그
ShouldBindJSON()	JSON 본문	json
ShouldBindQuery()	쿼리 스트링	form
ShouldBindUri()	      URI 파라미터	uri
ShouldBindHeader()	HTTP 헤더	header
ShouldBind()	Content-Type 기반 자동 선택	form, json
```

### 유효성 검사 

https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme

Binding Tag
```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"required,gte=0,lte=150"`
    Password string `json:"password" binding:"required,min=8"`
}
```

유효성 검사 실패 시 에러 메시지 변경

```go
func formatValidationErrors(err error) map[string]string {
    errors := make(map[string]string)

    for _, e := range err.(validator.ValidationErrors) {
        field := e.Field()
        tag := e.Tag()

        switch tag {
        case "required":
            errors[field] = field + "은(는) 필수입니다"
        case "email":
            errors[field] = "유효한 이메일 주소를 입력하세요"
        case "min":
            errors[field] = field + "은(는) 최소 " + e.Param() + "자 이상이어야 합니다"
        case "max":
            errors[field] = field + "은(는) 최대 " + e.Param() + "자까지 가능합니다"
        case "eqfield":
            errors[field] = field + "이(가) 일치하지 않습니다"
        default:
            errors[field] = field + " 값이 유효하지 않습니다"
        }
    }

    return errors
}

r.POST("/signup", func(c *gin.Context) {
    var req SignupRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            c.JSON(http.StatusBadRequest, gin.H{
                "errors": formatValidationErrors(validationErrors),
            })
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"username": req.Username})
})
```

필수/ 옵션값에 대한 처리
```go
// 생성용 - 모든 필드 필수
type CreateProductRequest struct {
    Name        string  `json:"name" binding:"required"`
    Price       float64 `json:"price" binding:"required,gt=0"`
    Description string  `json:"description" binding:"required"`
}

// 수정용 - 포인터로 선택적 필드 구현
type UpdateProductRequest struct {
    Name        *string  `json:"name" binding:"omitempty,min=1"`
    Price       *float64 `json:"price" binding:"omitempty,gt=0"`
    Description *string  `json:"description"`
}
```

커스텀 유효성 검사기
[binding.Validator.Engine()]
```go
// 한국 휴대폰 번호 검증
var phoneRegex = regexp.MustCompile(`^01[0-9]-?\d{3,4}-?\d{4}$`)

func phoneValidator(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    return phoneRegex.MatchString(phone)
}

type CreateUserRequest struct {
Name  string `json:"name" binding:"required"`
Phone string `json:"phone" binding:"required,phone"`
}

func main() {
r := gin.Default()

// 커스텀 검증기 등록
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
v.RegisterValidation("phone", phoneValidator)
}

r.POST("/users", func(c *gin.Context) {
var req CreateUserRequest

if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusCreated, gin.H{
"name":  req.Name,
"phone": req.Phone,
})
})

r.Run(":8080")
}
```

파라미터 포함 커스텀 검증기
```go
func startsWithValidator(fl validator.FieldLevel) bool {
    param := fl.Param()  // 태그에서 파라미터 가져오기
    value := fl.Field().String()
    return len(value) >= len(param) && value[:len(param)] == param
}

func main() {
r := gin.Default()

if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
v.RegisterValidation("startswith", startsWithValidator)
}

// 사용 예
type Request struct {
Code string `json:"code" binding:"required,startswith=PRD"`
}

}
```
구조체 필드 간 관계 검증
```go
type DateRangeRequest struct {
    StartDate string `json:"start_date" binding:"required"`
    EndDate   string `json:"end_date" binding:"required"`
}

func dateRangeValidator(sl validator.StructLevel) {
    req := sl.Current().Interface().(DateRangeRequest)

    // 날짜 파싱 (간략화)
    if req.StartDate > req.EndDate {
        sl.ReportError(req.EndDate, "EndDate", "end_date", "daterange", "")
    }
}

func main() {
r := gin.Default()

if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
v.RegisterStructValidation(dateRangeValidator, DateRangeRequest{})
}

}
```

커스텀 바인더 구현
```
type customBinding struct{}

func (customBinding) Name() string {
return "custom"
}

func (customBinding) Bind(req *http.Request, obj any) error {
// 커스텀 바인딩 로직
return nil
}

// 사용
c.ShouldBindWith(&req, customBinding{})
```


### 파일 업로드 

파일 업로드 
```go
type UploadRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
}

const maxFileSize = 5 * 1024 * 1024 // 5MB

func getContentType(file *gin.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		var req UploadRequest

		// 폼 데이터 바인딩
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 다중 파일 가져오기
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "파일이 필요합니다"})
			return
		}

		files := form.File["files"]
		var savedFiles []gin.H

		for _, file := range files {
			// 파일 크기 검사
			if file.Size > maxFileSize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "파일 크기는 5MB를 초과할 수 없습니다"})
				return
			}

			// 확장자 검사
			ext := filepath.Ext(file.Filename)
			if !allowedExtensions[ext] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "허용되지 않는 파일 형식입니다"})
				return
			}

			// MIME 타입 검사
			contentType, err := getContentType(file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if contentType != "image/jpeg" && contentType != "image/png" && contentType != "application/pdf" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "이미지 또는 PDF 파일만 업로드 가능합니다"})
				return
			}

			// UUID로 고유한 파일명 생성
			newFilename := uuid.New().String() + ext
			dst := "./uploads/" + newFilename

			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			savedFiles = append(savedFiles, gin.H{
				"original_name": file.Filename,
				"saved_name":    newFilename,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"title":       req.Title,
			"description": req.Description,
			"files":       savedFiles,
		})
	})

	r.Run(":8080")
}
```

대용량 파일 스트리밍
```go
import (
    "io"
    "os"
)

r.POST("/stream-upload", func(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    defer file.Close()

    dst, err := os.Create("./uploads/" + header.Filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer dst.Close()

    // 스트리밍 복사
    written, err := io.Copy(dst, file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "filename": header.Filename,
        "size":     written,
    })
})
```



### 응답 처리

구조체 태그
```go
type User struct {
  ID        int    `json:"id"`
  Password  string `json:"-"`          // 응답에서 제외
  UpdatedAt string `json:"updated_at,omitempty"` // 빈 값(제로값)이면 제외
}
```


JSON 
```
c.JSON(http.StatusOK,gin.H{...}) // gin.H == map[string]any
c.JSON(http.StatusOK,struct) // 권장

users := []user{
    {ID: 1, Name: "A"},
    {ID: 2, Name: "B"},
    {ID: 3, Name: "C"}
}

c.JSON(http.StatusOK,users)

c.IndentedJSON(http.StatusOK, user) // 들여쓰기 된 JSON (개발 시 사용)

c.SecureJSON(http.StatusOK, []string{"A","B","C"}) // 배열 직접 반환 시 보안확보 (JSON 하이재킹 공격 방지) ??
c.PureJSON(http.StatusOK, gin.H{"html": "<b>Hello</b>"}) // HTML 특수문자 이스케이프

```

### 에러 처리


AbortWithStatusJSON (실패 시 return 도 필수)

  1. 에러 응답 보냄
  2. 이후 핸들러/미들웨어 실행 중단

에러 응답 구조를 작성하는것이 좋음 (Client 처리하기 쉬움)
```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

func respondError(c *gin.Context, status int, code, message string, details any) {
    c.AbortWithStatusJSON(status, ErrorResponse{
        Code:    code,
        Message: message,
        Details: details,
    })
}

r.GET("/users/:id", func(c *gin.Context) {
id := c.Param("id")

if id == "0" {
respondError(c, http.StatusNotFound, "USER_NOT_FOUND", "사용자를 찾을 수 없습니다", nil)
return
}

c.JSON(http.StatusOK, gin.H{"id": id})
})

```
