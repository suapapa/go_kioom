# Multi-stage build for kioom-mcp (SSE Transport mode)

# --- Build Stage ---
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

# git 및 SSL 인증서 빌드 타임 의존성 설치
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 의존성 파일을 먼저 복사하여 레이어 캐싱 활용
COPY go.mod go.sum ./
RUN go mod download

# 전체 소스코드 복사 (kioom-mcp는 루트 패키지 및 internal 패키지를 참조하므로 전체 컨텍스트 필요)
COPY . .

# Buildx가 자동으로 제공하는 빌드 타겟 아규먼트
ARG TARGETOS
ARG TARGETARCH

# kioom-mcp 바이너리 빌드 (정적 링크 및 디버그 정보 제거)
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /kioom-mcp ./cmd/kioom-mcp

# --- Runner Stage ---
FROM alpine:3.19

# Kiwoom API 호출은 HTTPS를 사용하므로 SSL 인증서 필수 설치
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 빌더 스테이지에서 컴파일 완료된 정적 바이너리 복사
COPY --from=builder /kioom-mcp /app/kioom-mcp

# 컨테이너 포트 노출 (SSE 기본 포트: 8765)
EXPOSE 8765

# NOTE: 런타임 자격 증명(KIOOM_APP_KEY, KIOOM_SECRET_KEY 등)은 보안을 위해 Dockerfile에
# ENV/ARG로 하드코딩하지 않습니다. 대신 다음과 같은 방법으로 주입하십시오:
#   1. 컨테이너 실행 시 환경 변수로 주입: -e KIOOM_APP_KEY="xxx" -e KIOOM_SECRET_KEY="xxx"
#   2. 컨테이너 실행 시 .env 파일 주입: --env-file .env
#   3. .env 파일을 컨테이너 /app/.env 경로에 볼륨 마운트: -v $(pwd)/.env:/app/.env

# 컨테이너 환경에서는 호스트에서 바인딩을 허용하기 위해 listen 주소를 0.0.0.0으로 해야 합니다.
ENTRYPOINT ["/app/kioom-mcp", "-transport", "sse", "-listen", "0.0.0.0:8765"]
