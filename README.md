Go 명령어
```bash
# 모듈 초기화
go mod init hello

# 프로그램 실행
go run main.go

# 빌드 (실행 파일 생성)
### GOOS=linux : 대상 운영체제  
### GOARCH=amd64 : 대상 아키텍처
go build -o hello main.go

# 코드 포맷팅
gofmt -w main.go

# 의존성 관리
    
    go get github.com/gin-gonic/gin # 최신 버전
    go get github.com/gin-gonic/gin@v1.9.0 # 특정 버전
    go get github.com/gin-gonic/gin@abc1234 # 특정 커밋
    go get github.com/gin-gonic/gin@v1 # 최신 마이너 버전
    go mod graph # 의존성 그래프
    go mod vendor # 의존성 vendor 폴더 복사 (오프라인 빌드, 의존성 고정)
    go mod download # 의존성 다운로드
    go mod tidy # 필요 의존성 추가, 불필요 의존성 제거
    export GOPRIVATE=github.com/mycompany/* # 비공개 저장소의 모듈을 사용
    replace github.com/username/mylib => ../mylib # 모듈 경로를 다른 경로로 대체 (로컬 개발)
  
# 작업 영역
go work init ./app ./lib
```


# 키워드

### 기초 문법
```
var
const
 - iota
scope
 - package , function, block

자료형 (Zero Value)
    int8 int16 int32 int64 int
    uint8 uint16 uint32 uint64 uint
    uintptr
    float32 float64 
    complex64 complex128
    bool 
    string
    byte rune

타입 변환
    int(x)
    strconv : 문자열 <-> 숫자 사이 변환
```

### 제어문
```
if else
switch case
    fallthrough
```

### 반복문

```
기본 : for i=0; i<5; i++{ ... }
while : for i<5 { ... }
무한루프 : for { ... }
배열, 슬라이스, 맵, 문자열 : for index, value := range array { ... }

break, continue, (label, goto)
```


### 함수

```
func 함수명(인자) (리턴) { 실행문 } 

defer : 리소스 정리
    후입 선출
슬라이스 전개
    numbers := []int{1,2,3,4,5,6}
    sum(numbers...)

함수 타입
    type Operation func(int, int) int
    func add(a, b int) int { return a + b }
    func sub(a, b int) int { return a - b } 
         
pass by value : 함수에 전달된 값은 복사된다.


클로저 (https://go.dev/blog/loopvar-preview)
    - 외부 스코프의 변수를 참조로 캡처 
    
    Go 1.22부터: 루프 변수가 반복마다 새로 생성됨
        변경 전: per-loop scope (루프 전체에서 하나의 변수)
        변경 후: per-iteration scope (반복마다 새 변수)
        
```


### 패키지

```
# 표준 패키지
    fmt, strings, strconv, time
    os, io, net/http, encoding/json, sync
# 노출
    대문자 : 공개
    소문자 : 비공개

# init
    패키지 정의 파일 순서대로 실행
    _ : 패키지 Import 시 init만 필요 시 (db driver)
```


### 복합 자료형 

```

slice(슬라이스)

    슬라이싱
    append(a, b...)
    copy(dst, src)
    append(s[:2], s[3:]...) # index 2 삭제
    다차원 : [][]int {{},{},{}}
    
    nil vs empty slice
        - 둘 다 len은 0
        - nil 의경우 (slice == nil true)
        - empty 의 경우 (s3 == nil false)
        - nil도 append 가능

Map(맵)

    초기화 방법
        make
            make(map[키 타입]값 타입)
        리터럴
            map[string]int{"a":1, "b":2}
        var
            var m map[string]int   # nil map
            
    사용법
        ages["Alice"]              # 읽기
        ages["Alice"] = "test"     # 쓰기
        ages["Alice"] = "replace"  # 수정
        age, ok := ages["Alice"]   # 키 존재 여부 확인 (ok - true/false) 
        delete (ages, "Alice")     # 삭제
        for name, age := range ages { ... } # 순회
            - 순서를 보장하지 않는다. [Key를 Slice에 넣고 정렬 및 순회 하면 됨]

        Value에 함수 타입
            operations := map[string]func(int, int) int{
                "add": func(a, b int) int { return a + b },
                "sub": func(a, b int) int { return a - b },
                "mul": func(a, b int) int { return a * b },
        }
                        
    주의
        nil Map 
            읽기 : 가능 (Zero Value)
            쓰기 : 패닉 (초기화해야 사용가능)
            
        Key Type : 비교 가능한 타입만 가능 - Slice, Map, Function 불가
        
        동시성 : 일반 Map을 여러 스레드가 사용하면 안됨
            mutex
            sync.Map
            
            
struct (구조체)

    선언
        type 구조체명 struct { 필드 }
    생성
        new
    포인터 접근
        p := Person{Name: "Alice", Age: 30}
        ptr := &p
        ptr.Name  // ((*ptr).Name과 같음) - Alice
    비교
       - == : 모든 피드가 같을 때 true
       - 비교 불가능 타입 존재 시 비교 불가
    
    임베딩
        - Composition
    태그
        메타 데이터 추가 (JSON 직렬화)
            `json:"name"` : 필드 명
            `json:",omitempty" : 제로값이면 생략
            `json:"-" :json에서 제외
            
    익명 구조체 (테스트 코드 작성 시 유용)
        struct { X, Y int } {10,20}

Method (메서드)
    타입에 행동을 추가 (특정 타입에 연결된 함수)     
    
    선언
        func (타입) 함수명(파라미터) (반환) { 구현 }
        
        값 리시버 : 값의 복사본 (원본 수정 불가)
            func (p Person) SetAge(age int) { age int }
        포인터 리시버 : 원본 수정
             func (p *Person) SetAge(age int) { age int }
             
        선택 기준 - 얌전히 포인터 쓰세요
            상태 변경 : 포인터
            구조체가 클 때 : 포인터 (복사비용 감소)
            일관성 (다른 메서드도 다 포인터일때)     
            읽기 : 값
    특징     
        리시버 타입에 따라 실 사용 시 구조체를 자동 변환하여 사용한다 (자동 역참조)
            - 메서드 정의 시 포인터 리시버 / 값 리시버를 신중하게 결정
            - 사용 할 땐 구조체 선언하고 편안하게 호출!!
                c := Counter{} // 구조체 초기화 ( &Counter{} 도 상관 X)
                c.Increment()  // 그냥 호출 - 리시버 타입에 따라 참조 넘길지 값 넘길지 결정되어 있음!! 
    
        포인터 리시버는 nil일 수 있다.
            메서드 내에서 struct에 대한 nil 체크를 수행할 수 있음
                func (l *IntList) Sum() int {
                    if l == nil {
                        return 0
                    }
                    return l.Value + l.Next.Sum()
                }
                
        구조체 뿐 아니라 어떤 타입에도 메서드 추가 가능
            같은 패키지 내 정의된 타입에만 메서드를 추가할 수 있다.'
            
        임베딩 시 내부 타입의 메서드도 승격된다.
        
        
        메서드 오버라이딩
            동일 이름 메서드 정의 시 덮어쓴다
```


### 인터페이스
```
암시적으로 구현 (메서드 시그니처만 맞으면 자동으로 인터페이스를 만족)

다형성 : 다양한 타입을 같은 방식사용
인터페이스 합성 : 인터페이스를 조합하여 새 인터페이스 생성
인터페이스의 값은 타입과 값으로 구성됨
    인터페이스에 nil 포인터 할당 시 인터페이스 자체는 nil이 아님
    
작은 인터페이스를 조합하여 사용하는것을 권장


사용법
    Stringer 인터페이스  : 출력 형식을 제어
        type Person struct {
            Name string
            Age  int
        }
        
        func (p Person) String() string {
        return fmt.Sprintf("%s (%d세)", p.Name, p.Age)
        }

        p := Person{Name: "Alice", Age: 30}
        fmt.Println(p)  // Alice (30세)

구체적 타입 추출
    func describe(s Speaker) {
        switch v := s.(type) {
        case Dog:
            fmt.Println("개:", v.Name)
        case Cat:
            fmt.Println("고양이:", v.Name)
        default:
            fmt.Println("알 수 없는 동물")
        }
    }
    
    
interface{} : 빈 인터페이스

    모든 타입은 빈 인터페이스를 만족한다. - 어떤 타입의 값이든 저장 가능
    any == interface{}
    동적 데이터 다룰 때 유용
    
    주의점
        타입 안정성 손실 (컴파일 시 타입 검사 이점을 잃음)
            -> 1.18이후 제네릭 지원 : 우선 검토 고려
            
    사용 사례        
        사용 시 타입 단언 필요
            func double(v any) any {
                switch x := v.(type) {
                case int:
                    return x * 2
                case string:
                    return x + x
                default:
                    return nil
                }
            }    
        
        JSON 처리 (encofing/json)
            jsonStr := `{"name": "Alice", "age": 30, "active": true}`
        
            var data map[string]any
            json.Unmarshal([]byte(jsonStr), &data)
        
            fmt.Println(data["name"])    // Alice
            fmt.Println(data["age"])     // 30 (float64) # JSON의 숫자는 float64로 언마샬된다.
            fmt.Println(data["active"])  // true        
```
    
### 에러 처리

에러를 값으로 반환하여 에러처리를 명시적으로 만들었고 코드 흐름을 예측 가능하게 한다.

```
error 인터페이스를 구현함
    type error interface{ Error() string }
    
error를 값으로서 반환한다.
    - 함수 호출자가 처리한다.
    
사용
    생성
        errors.New() # 에러 생성
        fmt.Errorf("%s error", name) # 에러 포맷
        fmt.Errorf("%w", err) # 에러 래핑 
    확인
        errors.Is(err, os.ErrNotExist) # 에러 확인 (래핑된 에러에서도 찾는다)
    변환
        var pathError *PathError
        errors.As(err, &pathError) # err를 PathError 타입으로 변환하여 pathError에 복사한다.
    커스텀 에러
        type ValidationError struct {
            Field   string
            Message string
        }
        func (e ValidationError) Error() string {
            return fmt.Sprintf("%s: %s", e.Field, e.Message)
        }

    센티널 에러 : 패키지 수준에서 정의된 에러
        var (
            ErrNotFound = errors.New("not found")
            ErrForbidden = errors.New("forbidden")
        )    
        func getUser(id int) (string, error) {
            if id == 0 {
                return "", ErrNotFound
            }
            if id < 0 {
                return "", ErrForbidden
            }
            return "Alice", nil
        }

에러 처리 패턴
    1. 에러 발생 시 반환
        err != nil {
            return err
        }
    2. 에러 무시할 땐 명시한다.( 꼭 필요 시 ) 
        _ = logger.Close()           
```

### 패닉과 복구

    프로그램 정상 실행 불가 시 발생 Recover를 통해 프로그램 크래시를 방지한다.

    예시
        배열/슬라이스 범위 초과
        nil 포인터 역참조 (nil 포인터 참조 시)
        타입 단언 실패
        닫힌 채널에 송신
        nil Map에 쓰기

    Recover (복구)
        패닉을 잡아 프로그램 종료를 방지한다. (defer 내에서 동작)

            func riskyOperation() {
                defer func() {
                    if r := recover(); r != nil {
                        fmt.Println("에러 발생, 복구 중:", r)
                    }
                    fmt.Println("정리 작업 수행")
                }()
            
                fmt.Println("위험한 작업 시작")
                panic("예상치 못한 에러")
            }
            
            func main() {
                riskyOperation()
                fmt.Println("프로그램 계속 실행")
            }


    사용 기준
        1. 프로그램 초기화 실패 시
        2. 프로그래머 실수
        3. Must접두사 함수 등의 컨벤션 지정

    절대 사용 금지
        1. 예상 가능한 에러에 panic
        2. Validation 체크 
        3. 외부 시스템 오류
        4. 라이브러리 패닉 -> 에러로 변환하여 처리하여 비정상 종료를 방지 추천
          