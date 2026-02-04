

# 댓글
### 댓글 생성
curl -X POST http://localhost:8080/api/v1/posts/1/comments \
-H "Content-Type: application/json" \
-d '{
"content": "좋은 글이네요!",
"author": "독자1"
}'

### 댓글 목록 조회
curl http://localhost:8080/api/v1/posts/1/comments

### 댓글 수정
curl -X PUT http://localhost:8080/api/v1/posts/1/comments/1 \
-H "Content-Type: application/json" \
-d '{"content": "수정된 댓글입니다"}'

### 댓글 삭제
curl -X DELETE http://localhost:8080/api/v1/posts/1/comments/1

# 댓글 목록 조회 (대댓글 포함)
curl http://localhost:8080/api/v1/posts/1/comments
