# REST

- URL로 리소스를 식별하고, HTTP 메서드로 행위를 표현

- URL : 명사(복수형) / 행위 : 메소드
- 리소스 간 관계는 URL 경로로 표현



### 대댓글 구현

```
인접 목록 (ParentID)	구현 단순, 삽입 빠름	조회 시 재귀/반복 필요
경로 열거 (Path)	조회 빠름	삽입/수정 복잡
중첩 집합	범위 조회 빠름	삽입 매우 복잡
```


# 통합 테스트

### 테스트 데이터베이스 생성
psql -U postgres -c "CREATE DATABASE godb_test OWNER gouser;"

### 테스트 실행
go test -v ./internal/handler/...

### 테스트 커버리지
go test -cover ./...



##### 정리
게시글 목록	GET /api/v1/posts
게시글 생성	POST /api/v1/posts
게시글 조회	GET /api/v1/posts/:id
게시글 수정	PUT /api/v1/posts/:id
게시글 삭제	DELETE /api/v1/posts/:id
커서 페이징	GET /api/v1/posts/cursor
댓글 목록	GET /api/v1/posts/:postId/comments
댓글 생성	POST /api/v1/posts/:postId/comments
댓글 수정	PUT /api/v1/posts/:postId/comments/:id
댓글 삭제	DELETE /api/v1/posts/:postId/comments/:id