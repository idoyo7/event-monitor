# 시작 이미지로 공식 Go 이미지 사용
FROM golang:1.21.5 as builder

# 작업 디렉토리 설정
WORKDIR /app

# 프로젝트 소스 파일 복사
COPY . .

# Go 모듈 초기화 및 의존성 다운로드
# RUN go mod init event-monitoring && \
#     go mod tidy
RUN go mod tidy

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o event .

# 실행 이미지
FROM alpine:3.19.0  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 빌드된 실행 파일 복사
COPY --from=builder /app/event .

# 실행 명령어 설정
CMD ["./event"]

