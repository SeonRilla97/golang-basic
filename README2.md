# Go 언어 핵심 정리

## 명령어 요약

| 명령어 | 용도 |
|--------|------|
| `go mod init <name>` | 모듈 초기화 |
| `go run main.go` | 실행 |
| `go build -o <output> main.go` | 빌드 |
| `go mod tidy` | 의존성 정리 (추가/제거) |
| `go get <pkg>@<version>` | 패키지 설치 |
| `gofmt -w <file>` | 코드 포맷팅 |

---

## 변수와 타입

### 선언
```go
var x int           // Zero Value (0)
var x int = 10      // 명시적 초기화
x := 10             // 짧은 선언 (함수 내에서만)
const PI = 3.14     // 상수
```

### 기본 타입과 Zero Value
| 타입 | Zero Value |
|------|------------|
| `int`, `float64` | `0` |
| `bool` | `false` |
| `string` | `""` |
| `pointer`, `slice`, `map`, `chan` | `nil` |

### 타입 변환
- 명시적 변환만 허용: `int(x)`, `float64(y)`
- 문자열 변환: `strconv` 패키지 사용

---

## 제어문

### if / switch
```go
if x > 0 { }                    // 기본
if err := fn(); err != nil { }  // 초기화 포함

switch v {                      // fallthrough로 다음 case 실행
case 1: ...
default: ...
}
```

### 반복문 (for만 존재)
```go
for i := 0; i < 5; i++ { }      // 기본
for condition { }               // while처럼
for { }                         // 무한루프
for i, v := range slice { }    // 순회
```

---

## 함수

### 기본 형태
```go
func name(param int) (result int, err error) { }
```

### 핵심 개념
| 개념 | 설명 |
|------|------|
| **다중 반환** | `return value, err` |
| **defer** | 함수 종료 시 실행 (LIFO), 리소스 정리용 |
| **가변 인자** | `func sum(nums ...int)` / 호출: `sum(slice...)` |
| **Pass by Value** | 모든 인자는 복사됨 (포인터로 원본 수정) |
| **클로저** | 외부 변수 참조로 캡처 (Go 1.22+: 루프 변수 매 반복 새로 생성) |

---

## 패키지

### 가시성
- **대문자**: 공개 (Public)
- **소문자**: 비공개 (Private)

### init 함수
- 패키지 로드 시 자동 실행
- `import _ "pkg"`: init만 실행 (DB 드라이버 등록 등)

---

## 복합 자료형

### Slice
```go
s := []int{1, 2, 3}            // 생성
s = append(s, 4, 5)            // 추가
s = append(s[:i], s[i+1:]...)  // i번째 삭제
copy(dst, src)                 // 복사
```
> **nil slice vs empty slice**: 둘 다 `len=0`, nil만 `== nil` true

### Map
```go
m := make(map[string]int)      // 생성
m := map[string]int{"a": 1}    // 리터럴
v, ok := m["key"]              // 존재 확인
delete(m, "key")               // 삭제
```
> **주의**: nil map 읽기 가능, 쓰기 panic / 순회 순서 미보장 / 동시성: `sync.Map` 또는 mutex

### Struct
```go
type Person struct {
    Name string `json:"name,omitempty"`  // 태그
    Age  int    `json:"-"`               // JSON 제외
}
p := Person{Name: "Alice"}
ptr := &p
```
> **임베딩**: 상속 대신 조합 (Composition)

---

## 메서드

### 리시버 타입
```go
func (p Person) Read() string { }   // 값 리시버: 복사본
func (p *Person) Update() { }       // 포인터 리시버: 원본 수정
```

### 선택 기준
| 상황 | 리시버 |
|------|--------|
| 상태 변경 | 포인터 |
| 큰 구조체 | 포인터 (복사 비용) |
| 읽기 전용 | 값 |
| **일관성** | 하나로 통일 권장 |

> **자동 역참조**: `c.Method()` 호출 시 리시버 타입에 맞게 자동 변환

---

## 인터페이스

### 핵심 특징
- **암시적 구현**: 메서드 시그니처만 맞으면 자동 만족
- **작은 인터페이스** 조합 권장 (io.Reader, io.Writer 등)
- **`any` = `interface{}`**: 모든 타입 수용 (타입 안정성 손실)

### 타입 단언 / 스위치
```go
v, ok := i.(ConcreteType)      // 타입 단언
switch v := i.(type) { }       // 타입 스위치
```

### 자주 쓰는 인터페이스
- `Stringer`: `String() string` - fmt 출력 형식 제어
- `error`: `Error() string`

---

## 에러 처리

### 원칙
- **에러는 값**: 반환하여 호출자가 처리
- **예외 없음**: panic은 진짜 비정상 상황에만

### 사용법
```go
errors.New("message")              // 생성
fmt.Errorf("context: %w", err)     // 래핑
errors.Is(err, target)             // 비교 (래핑 체인 탐색)
errors.As(err, &targetType)        // 타입 변환
```

### 센티널 에러
```go
var ErrNotFound = errors.New("not found")  // 패키지 수준 정의
```

---

## 패닉과 복구

### Panic 발생 상황
- 배열 범위 초과
- nil 포인터 역참조
- nil map 쓰기
- 닫힌 채널 송신

### Recover
```go
defer func() {
    if r := recover(); r != nil {
        // 복구 처리
    }
}()
```

### 사용 원칙
| O 사용 | X 금지 |
|--------|--------|
| 초기화 실패 | 예상 가능한 에러 |
| 프로그래머 실수 | Validation 실패 |
| Must* 함수 컨벤션 | 외부 시스템 오류 |

---

## 요약 다이어그램

```
┌─────────────────────────────────────────────────────────────┐
│                        Go 핵심 구조                          │
├─────────────────────────────────────────────────────────────┤
│  타입 시스템                                                 │
│  ├─ 기본 타입 (int, string, bool...)                        │
│  ├─ 복합 타입 (slice, map, struct)                          │
│  └─ 인터페이스 (암시적 구현, 다형성)                          │
├─────────────────────────────────────────────────────────────┤
│  함수와 메서드                                               │
│  ├─ 다중 반환, defer, 클로저                                 │
│  └─ 값/포인터 리시버 → 자동 역참조                           │
├─────────────────────────────────────────────────────────────┤
│  에러 처리                                                   │
│  ├─ error는 값 → 명시적 처리                                 │
│  └─ panic/recover → 진짜 비정상만                            │
└─────────────────────────────────────────────────────────────┘
```