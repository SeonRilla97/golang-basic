어플리케이션의 에러처리

Logging

    1. AppError을 정의하여 비즈니스 로직에 적절한 로직을 작성
    2. Client에 내려줄 HTTP 전용 Error struct를 정의
    3. Middle Ware 를 통해 AppError -> HTTP Error 로 변환하는 과정을 거친다.

Recovery

    1. 고루틴은 에러 핸들링 미들웨어를 거치지 않는다 -> 별도 Recovery