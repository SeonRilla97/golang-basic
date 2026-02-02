# GORM

객체와 관계형 데이터베이스를 연결 (객체를 데이터베이스 테이블의 행(row)으로 자동 변환)

DB 생성 (Docker)

```go
docker run --name go-postgres \
  -e POSTGRES_USER=gouser \
  -e POSTGRES_PASSWORD=gopassword \
  -e POSTGRES_DB=godb \
  -p 5432:5432 \
  -d postgres:16


docker exec -it go-postgres psql -U gouser -d godb -c "SELECT version();"

  
```


DB 연결

    - 접속 정보 환경변수로 관리
    - 연결 코드를 별도 함수로 분리
    - DB 설정 가능
    - Connection Pool 설정

```go
// filename: database/database.go
package database

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_USER", "gouser"),
        getEnv("DB_PASSWORD", "gopassword"),
        getEnv("DB_NAME", "godb"),
        getEnv("DB_PORT", "5432"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// 느린 쿼리 로깅 (200ms 이상)
		Logger: logger.Default.LogMode(logger.Info),

		// 테이블명 단수 사용 (users -> user)
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},

		// 트랜잭션 비활성화 (성능 향상)
		SkipDefaultTransaction: true,

		// 생성 시 기본값 사용
		PrepareStmt: true,
	})
    if err != nil {
        log.Fatal("데이터베이스 연결 실패:", err)
    }

    sqlDB, _ := DB.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    log.Println("데이터베이스 연결 성공")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

DB 모델 정의

```go

// 자주 사용하는 필드들이 정의되어 있음 (ID, CreatedAt, UpdatedAt, DeletedAt) & Soft Delete 지원
gorm.Model

// 구조체 정의

primaryKey	기본키 지정	gorm:"primaryKey"
column	컬럼명 지정	gorm:"column:user_name"
type	컬럼 타입 지정	gorm:"type:varchar(100)"
size	문자열 크기	gorm:"size:100"
default	기본값	gorm:"default:0"
not null	NULL 불허	gorm:"not null"
index	인덱스 생성	gorm:"index"
uniqueIndex	유니크 인덱스	gorm:"uniqueIndex"
index:idx_name	인덱스 이름 지정	gorm:"index:idx_email"


// 테이블 정의
GORM은 구조체 이름을 복수형, snake_case로 변환하여 테이블 이름을 만듭니다.
    테이블 이름을 직접 지정하고 싶다면 TableName() 메서드를 정의


// 컬럼 정의
snake_case로 변환
    컬럼 이름을 직접 지정하려면 column 태그를 사용

크기를 지정하려면 size 또는 type 태그
데이터베이스에 저장하지 않을 필드는 - 태그



// 필드 타입 매핑
int, int64	bigint
int32	integer
float64	double precision
string	varchar(255)
bool	boolean
time.Time	timestamp
[]byte	bytea
```

### 마이그레이션


##### autoMigrate : 개발환경에서만 사용
```go
구조체를 기반으로 테이블을 자동 생성하거나 수정 (컬럼 삭제는 발생하지 않음)
package database

import (
    "log"

    "myapp/models"
)
func Migrate() {

db.AutoMigrate(&User{}, &Post{})
log.Println("개발 환경: AutoMigrate 실행")
} else {
log.Println("프로덕션 환경: AutoMigrate 건너뜀")
}
}



func Migrate() {
    if os.Getenv("ENV") == "development" {
        err := DB.AutoMigrate(
            &models.User{},
            &models.Post{},
            &models.Comment{},
        )
        if err != nil {
			log.Fatal("마이그레이션 실패:", err)
        }
        log.Println("마이그레이션 완료")
    }
}

```

세밀한 제어 db.Migrator()
```go
m := db.Migrator()

m.CreateTable(&User{}) // 테이블 생성
m.DropTable(&User{}) // 테이블 삭제
m.RenameTable(&User{}, "members") // 테이블 이름 변경
m.AddColumn(&User{}, "Bio") // 컬럼 추가
m.DropColumn(&User{}, "Bio") // 컬럼 삭제
m.RenameColumn(&User{}, "Name", "FullName") // 컬럼 이름 변경
m.CreateIndex(&User{}, "idx_email") // 인덱스 생성
m.DropIndex(&User{}, "idx_email") // 인덱스 삭제

```

##### golang-migrate

```go
# golang-migrate 설치
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 마이그레이션 파일 생성
migrate create -ext sql -dir migrations -seq add_users_table

# 마이그레이션 실행
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up


-- migrations/001_create_users.sql
CREATE TABLE users (
id SERIAL PRIMARY KEY,
name VARCHAR(100) NOT NULL,
email VARCHAR(255) UNIQUE NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/002_add_age_to_users.sql
ALTER TABLE users ADD COLUMN age INTEGER DEFAULT 0;


```

### CRUD

##### Create

```
db.Create(&user) // 단건
db.Select("Name", "Email").Create(&user) // 특정 필드만
result := db.Create(&users) // 다건 (구조체 슬라이스 전달)
db.CreateInBatches(users, 100) // 100개씩 배치로 생성 (보통 100~ 1000 사이로 제어 - 메모리/성능 균형 맞춤)

Upsert
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},           // 충돌 감지 컬럼
    DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),  // 업데이트할 컬럼
}).Create(&user)

// 출동 시 아무것도 안할 때
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

```

- default Value - 필드에 default 태그
- 에러 처리 시 String포함 여부 확인
    if strings.Contains(result.Error.Error(), "duplicate key") {

##### Read
```go
# 단건 조회
db.First(&user, 1) // 기본키로 조회 (PK 정렬, 1개 데이터)
db.Last(&user) // 기본키로 조회 (PK 정렬, 마지막 데이터, 1개 데이터)
db.Take(&user) // 정렬 없이 조회  


# 다건 조회
db.Find(&users) // 데이터 없으면 빈 슬라이스 반환 (에러 안남)
db.Find(&users, []int{1, 2, 3})  // SELECT * FROM users WHERE id IN (1, 2, 3)

# 조건 조회
db.Where("age > ?", 25).Find(&users) // 문자열 조건
db.Where("age > ? AND name LIKE ?", 25, "%김%").Find(&users) // AND 조건
db.Where(&User{Name: "홍길동", Age: 30}).Find(&users) // 구조체 조건 (제로값은 무시됨)
db.Where(map[string]interface{}{"name": "홍길동", "age": 0}).Find(&users)// 맵 조건 (제로값도 포함)
db.Where("age > ?", 30).Or("name = ?", "홍길동").Find(&users) // Or 조건 추가
db.Not("age > ?", 30).Find(&users) // Not 조건
db.Select("name", "email").Find(&users) // 특정 필드만 조회
db.Select("name", "age * 2 as double_age").Find(&users) // 별칭 사용
db.Model(&User{}).Distinct("age").Pluck("age", &ages) // 중복 제거

# 정렬과 페이징
db.Order("age desc, name").Find(&users) // 다중정렬, 오름차순/내림차순
db.Offset(10).Limit(10).Find(&users) // 페이징

# 집계 함수
db.Model(&User{}).Where("age > ?", 25).Count(&count)  // 개수
db.Model(&User{}).Pluck("name", &names) // 단일 컬럼 값을 슬라이스로 가져옴
db.Model(&User{}).Select("age, count(*) as count").Group("age").Having("count(*) > ?", 1).Scan(&results) // Group & Having

# 스캔
db.Model(&User{}).Select("name", "email").First(&dto) // 조회 결과를 다른 구조체에 매핑
row := db.Model(&User{}).Where("id = ?", 1).Select("name", "age").Row() // 단일 행 직접 스캔
rows, _ := db.Model(&User{}).Select("name", "age").Rows() // 다중 행 직접 스캔

# 쿼리 디버깅
db.Debug().Where("age > ?", 25).Find(&users)

// 레코드가 없을 때 발생 에러 -> 에러 핸들링 시 필요
gorm.ErrRecordNotFound
```



##### Update
```go
db.Save(&user) //모든 필드 업데이트, 기본키 존재하면 UPDATE, 없으면 INSERT
db.Model(&user).Update("age", 32) // 단일 필드 업데이트
db.Model(&user).Updates(User{Name: "업데이트", Age: 33}) // 구조체가 ZeroValue(0,"",false) 면 무시
db.Model(&user).Updates(map[string]interface{}{"name": "테스트","age":  0,})  // Zero Value 업데이트 필요 시
db.Model(&user).Select("name").Updates(User{Name: "선택", Age: 99}) // 특정 필드만 업데이트
db.Model(&user).Omit("age").Updates(User{Name: "제외", Age: 99}) // 특정 필드 제외 (age)
db.Model(&User{}).Where("age >= ?", 30).Update("status", "senior") // 다건 업데이트
db.Model(&user).Update("age", gorm.Expr("age + ?", 1)) // SQL 표현식으로 업데이트
db.Model(&user).UpdateColumn("name", "빠름") // Hook 건너 뜀, UpdatedAt 갱신 X
if tx.Statement.Changed("Name") {fmt.Println("이름이 변경되었습니다")} // 업데이트 전 변경 확인
```


##### DELETE

```go
db.Delete(&user) // gorm.Model 사용 시 Soft Delete -> deleted_at 삭제 시간 기록
db.Delete(&User{}, []int{1, 2, 3}) // 기본키 직접 지정하여 삭제, [슬라이스, 기본형 가능]
result := db.Where("age > ?", 100).Delete(&User{}) // 조건 지정하여 삭제
db.Where("1 = 1").Delete(&User{}) // 테이블 데이터 전체 삭제 시 명시하지 않으면 에러 발생
db.Unscoped().Delete(&user) // Hard Delete 
db.Select("Posts").Delete(&user) // user와 연관된 Posts도 삭제
```
Gorm DeletedAt 필드 커스터마이징
```
type User struct {
    ID        uint
    Name      string
    DeletedAt gorm.DeletedAt `gorm:"index"`  // 기본
}

// 또는 soft delete 플래그 사용
type User struct {
    ID        uint
    Name      string
    IsDeleted bool `gorm:"default:false"`
}

// Soft Delete 비활성화
type User struct {
    ID   uint
    Name string
    // DeletedAt 필드 없음 = Soft Delete 없음
}
```

삭제 복구

```go
func RestoreUser(db *gorm.DB, userID uint) error {
    result := db.Model(&User{}).
        Unscoped().
        Where("id = ?", userID).
        Update("deleted_at", nil)

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return errors.New("복구할 사용자를 찾을 수 없습니다")
    }

    return nil
}
```



### 연관 관계

##### 1:N (일대다)

HasMany
    - 1에서의 선언 : Slice

BelongsTo
    - N에서의 선언 : 외래키 필드 및 참조 필드를 선언

다른 필드를 참조 :  references 태그를 사용

외래키 자동 추론
    - [모델명]ID
    - column : 외래키 컬럼명 지정
    - foreignKey : 왜래키 필드 지정

외래키 제약조건
    - OnUpdate:CASCADE	부모 레코드 수정 시 함께 수정
    - OnDelete:CASCADE	부모 레코드 삭제 시 함께 삭제
    - OnDelete:SET NULL	부모 레코드 삭제 시 NULL로 설정
    - OnDelete:RESTRICT	자식 레코드가 있으면 삭제 불가

연관 데이터 생성 방법
```go
db.Model(&user).Association("Posts").Append(&post) // append
db.Create(&Post{Title: "또 다른 글", UserID: user.ID}) // Create

```


##### 다대다

many2many 사용하여 다대다 가능하고, Join Table 자동 생성되나, 관리되지 못하기에 수동생성

```go
// 조인 테이블 모델
type PostTag struct {
    PostID    uint      `gorm:"primaryKey"`
    TagID     uint      `gorm:"primaryKey"`
    CreatedAt time.Time // 추가 필드
    CreatedBy string    // 추가 필드
}

type Post struct {
    gorm.Model
    Title string
    Tags  []Tag `gorm:"many2many:post_tags"`
}

type Tag struct {
    gorm.Model
    Name string
}

// SetupJoinTable로 조인 테이블 설정
db.SetupJoinTable(&Post{}, "Tags", &PostTag{})
db.Model(&post).Association("Tags").Append(&newTag) // 기존 게시글에 태그 연결

// 기존 태그 찾아 추가
var tag Tag
db.Where("name = ?", "Go").First(&tag)
db.Model(&post).Association("Tags").Append(&tag)


db.Model(&post).Association("Tags").Replace(&newTags) // 기존 태그를 모두 제거하고 새 태그로 교체
db.Model(&post).Association("Tags").Delete(&tag) // 연결 해제 (조인 테이블에서만 삭제, 태그 자체는 유지)
db.Model(&post).Association("Tags").Clear() // 모든 연관 해제
db.Preload("Tags").First(&post, 1) // 게시글에서 태그 조회
db.Preload("Posts").Where("name = ?", "Go").First(&tag) // 태그에서 게시글 조회


```


###### Preload

- 연관 데이터 조회 시 기본적으로 가져오지 않는다 (별도 요청 필요)
```go
db.Preload("Posts").Preload("Profile").First(&user, 1) //단건 조회 (Posts, Profile 같이 조회) 
db.Preload("Posts").Find(&users) // 다건 조회
db.Preload("Posts.Comments").First(&user, 1) // 중첩 조회 (User -> Posts -> Comments 까지 조회)

// 연관 데이터 조회에 조건 걸기
// 최근 5개의 게시글만 로드
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
return db.Order("created_at DESC").Limit(5)
}).First(&user, 1)


db.Preload("Posts", "deleted_at IS NULL").First(&user, 1) // 삭제되지 않은 게시글만 로드
db.Preload("Posts", "views > ?", 100).First(&user, 1) // 특정 조건의 게시글만 로드


db.Preload(clause.Associations).First(&user, 1) //모든 연관 데이터 로드

Join
    -> Joins 일대일 관계에서만 사용 가능
// Preload: 2개의 쿼리 실행
// SELECT * FROM users WHERE id = 1
// SELECT * FROM posts WHERE user_id = 1
db.Preload("Posts").First(&user, 1)

// Joins: 1개의 쿼리 실행 (단, 일대일 관계만)
// SELECT users.*, profiles.* FROM users LEFT JOIN profiles ON ...
db.Joins("Profile").First(&user, 1)
```
 


N+1 는 Preloads 를 이용하여 해결한다.
연관 데이터 조회 시 N+1 문제를 해결하기 위함
```go
var users []User
db.Preload("Posts").Find(&users)
```


사용 예시 
```go
func GetPost(db *gorm.DB, postID uint) (*Post, error) {
    var post Post

    err := db.
        Preload("User").                           // 작성자
        Preload("Tags").                           // 태그
        Preload("Comments", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at DESC").Limit(10)  // 최근 댓글 10개
        }).
        Preload("Comments.User").                  // 댓글 작성자
        First(&post, postID).Error

    if err != nil {
        return nil, err
    }

    return &post, nil
}
다음 단계
```
##### Association

연관 데이터 직접 관리 (CRUD)

- 다대다 관계에서는 조인 테이블에 대한 제어를 함


Association.Find	연관 데이터 조회
Association.Append	연관 추가 (기존 연관 유지, 신규 추가)
Association.Replace	연관 교체 (기존 연관 모두 제거, 새로 설정)
Association.Delete	특정 연관 제거 (지정된 연관만 제거)
Association.Clear	모든 연관 제거 (관련 연관 모두 제거)

db.Model(&user).Association("Posts").Unscoped().Find(&posts) // Soft Delete 연관도 포함하여 제어

### 고급 쿼리



##### Raw SQL

조인,서브쿼리 / DB별 특화 기능/ 성능 최적화 튜닝 / 레거시 쿼리 마이그레이션 필요시

```go
## Raw Query (Raw SQL을 사용할 때는 항상 파라미터 바인딩(?)을 사용 - SQL Injection 대응)
db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&users)
db.Raw("SELECT * FROM users WHERE age > @age",sql.Named("age", 25),).Scan(&users)
db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "새사용자", 30)

## Row & Rows
row := db.Raw("SELECT name, age FROM users WHERE id = ?", 1).Row() // 단일 행 
rows, _ := db.Raw("SELECT name, age FROM users WHERE age > ?", 20).Rows() // 다중 행

## From / Where 절에서의 서브 쿼리

// WHERE
subQuery := db.Model(&User{}).Select("AVG(age)")
var users []User
db.Where("age > (?)", subQuery).Find(&users)

// FROM
subQuery := db.Table("users").Select("name, age")
db.Table("(?) as u", subQuery).Where("u.age > ?", 25).Find(&results)


// 테스트(??)

1. SQL 확인
stmt := db.Session(&gorm.Session{DryRun: true}).
Where("age > ?", 25).
Find(&users).
Statement

sql := stmt.SQL.String()
vars := stmt.Vars

2. SQL 문자열 생성
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
return tx.Model(&User{}).Where("id = ?", 1).Find(&users)
})

fmt.Println(sql)
// SELECT * FROM "users" WHERE id = 1 AND "users"."deleted_at" IS NULL
```


##### 트랜잭션

1. Create , Update Delete에 자동 트랜잭션을 적용함 (GORM)
2. 자동 트랜잭션 비활성화 시 성능 향상


- Transaction 메소드 사용
  - nil : Commit
  - err : Rollback

- SavePoint
  - 중간 저장점을 통한 부분 롤백


트랜잭션 격리 수준

```go
LevelDefault	데이터베이스 기본값
LevelReadUncommitted	커밋되지 않은 데이터 읽기 가능
LevelReadCommitted	커밋된 데이터만 읽기
LevelRepeatableRead	트랜잭션 내 반복 읽기 일관성
LevelSerializable	완전 직렬화
```

DB 락
```go
// 배타적 잠금 (FOR UPDATE)
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)

// 공유 잠금 (FOR SHARE)
db.Clauses(clause.Locking{Strength: "SHARE"}).Find(&users)

// 대기 없음 (NOWAIT)
db.Clauses(clause.Locking{
    Strength: "UPDATE",
    Options:  "NOWAIT",
}).Find(&users)
```


##### Scope

자주 사용하는 쿼리 로직을 함수화 한다. 

페이징, 정렬, 검색, 필터, 비즈니스 관련 스코프(모델)
```go
func GetUsersHandler(c *gin.Context) {
    // 쿼리 파라미터 파싱
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    keyword := c.Query("q")
    role := c.Query("role")

    // 필터 생성
    filter := UserFilter{
        Name: keyword,
        Role: role,
    }

    // 스코프 조합하여 조회
    var users []User
    var total int64

    db.Model(&User{}).
        Scopes(FilterUsers(filter)).
        Count(&total)

    db.Scopes(
        FilterUsers(filter),
        Paginate(page, pageSize),
        OrderBy("created_at", "desc"),
    ).Find(&users)

    c.JSON(200, gin.H{
        "data":       users,
        "total":      total,
        "page":       page,
        "page_size":  pageSize,
    })
}
```


##### HOOK

모델의 생명주기에 개입 (콜백함수) 특정 시점에 실행될 함수

비밀번호 해싱, 타임스탬프 갱신, 유효성 검사, UUID, 슬러그, 연관데이터 처리, 조회 후 데이터 가공, Soft Delete 전처리, 감사로그

```go
Hook	실행 시점
BeforeCreate	Create 전
AfterCreate	Create 후
BeforeUpdate	Update 전
AfterUpdate	Update 후
BeforeSave	Create/Update 전
AfterSave	Create/Update 후
BeforeDelete	Delete 전
AfterDelete	Delete 후
AfterFind	조회 후

## Hook 건너뛰기
db.Model(&user).UpdateColumn("name", "새이름") // Hook 건너뛰기
db.Session(&gorm.Session{SkipHooks: true}).Create(&user) // 세션으로 비활성화

```

