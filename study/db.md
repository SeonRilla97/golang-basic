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



// 레코드가 없을 때 발생 에러 -> 에러 핸들링 시 필요
gorm.ErrRecordNotFound
```