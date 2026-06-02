# `kioom-mcp` — Kiwoom MCP server

[`go_kioom`](https://github.com/suapapa/go_kioom) Kiwoom REST 클라이언트를 **Model Context Protocol(MCP)** 서버로 노출합니다. 도구 정의는 `internal/mcpkioom`에 있으며, **stdio**와 **HTTP + SSE** 두 가지 전송을 같은 서버 인스턴스로 제공합니다.

## 요약

| 전송        | 용도 |
|------------|------|
| `stdio` (기본) | Cursor, Claude Desktop 등 로컬 MCP 클라이언트 |
| `sse`      | HTTP에서 SSE 세션을 열고 싶을 때 (원격/브라우저 클라이언트 등) |

자격 증명과 검증 규칙은 [`kioom-cli`](../kioom-cli/README.md)와 동일한 환경 변수를 사용합니다.

## 설치

```bash
go install github.com/suapapa/go_kioom/cmd/kioom-mcp@latest
```

## 빌드

```bash
go build -o bin/kioom-mcp ./cmd/kioom-mcp
```

## 환경 변수

| 변수 | 설명 |
|------|------|
| `KIOOM_APP_KEY` | Kiwoom 앱 키 (**필수**) |
| `KIOOM_SECRET_KEY` | 시크릿 키 (**필수**) |
| `KIOOM_TOKEN` | 이미 발급된 Bearer 토큰이 있으면 설정 (선택) |
| `KIOOM_MOCK` | `true`이면 모의투자(mock) API 도메인 사용 |
| `KIOOM_MCP_SSE_TOKEN` | SSE 클라이언트 인증에 사용할 Bearer 토큰 (선택) |

`KIOOM_APP_KEY`와 `KIOOM_SECRET_KEY`가 없으면 프로세스는 즉시 종료됩니다.

## 플래그

| 플래그 | 기본값 | 설명 |
|--------|--------|------|
| `-transport` | `stdio` | `stdio` 또는 `sse` |
| `-listen` | `127.0.0.1:8765` | `-transport=sse`일 때 `host:port` |
| `-sse-path` | `/` | SSE 핸들러를 붙일 URL 경로. 루트가 아니면 내부에서 끝에 `/`가 붙을 수 있음(Go 1.22+ `ServeMux` 서브트리 매칭) |
| `-sse-token` | (환경변수값) | SSE 클라이언트 인증용 Bearer 토큰. 설정 시 클라이언트는 HTTP Header (`Authorization: Bearer <토큰>`) 또는 URL 쿼리 파라미터 (`?token=<토큰>` 또는 `?auth=<토큰>`)를 제공해야 합니다. |

SIGINT/SIGTERM으로 종료하면 **SSE 모드**에서는 `http.Server.Shutdown`으로 정리합니다(최대 10초).

## 실행 예

### stdio (로컬 MCP)

```bash
export KIOOM_APP_KEY="..."
export KIOOM_SECRET_KEY="..."
# 선택: export KIOOM_TOKEN="..." 
# 선택: export KIOOM_MOCK=true

kioom-mcp
# 또는 명시적으로
kioom-mcp -transport stdio
```

Claude Desktop 등 MCP 설정 예시(개념):

```json
{
  "mcpServers": {
    "kioom": {
      "command": "kioom-mcp",
      "env": {
        "KIOOM_APP_KEY": "your-app-key",
        "KIOOM_SECRET_KEY": "your-secret-key",
        "KIOOM_MOCK": "true"
      }
    }
  }
}
```

경로는 클라이언트 버전에 따라 다를 수 있으니 해당 앱 문서를 따릅니다.

### SSE (HTTP)

```bash
export KIOOM_APP_KEY="..."
export KIOOM_SECRET_KEY="..."

# 인증 토큰 없이 SSE 구동
kioom-mcp -transport sse -listen 127.0.0.1:8765 -sse-path /

# 인증 토큰을 지정하여 SSE 구동
kioom-mcp -transport sse -listen 127.0.0.1:8765 -sse-path / -sse-token "my-secret-token"
```

### Docker를 사용한 SSE 구동

제공되는 Dockerfile을 사용하여 `kioom-mcp`를 컨테이너 환경에서 빌드하고 구동할 수 있습니다. 

#### 1. 컨테이너 이미지 빌드
컨테이너 빌드는 프로젝트 **루트 디렉터리**에서 실행해야 합니다.
```bash
docker build -t kioom-mcp-sse -f cmd/kioom-mcp/Dockerfile .
```

#### 2. 컨테이너 실행
자격 증명 환경 변수를 넘겨서 로컬에서 빌드한 이미지를 실행합니다.
```bash
docker run -d \
  -p 8765:8765 \
  -e KIOOM_APP_KEY="your-app-key" \
  -e KIOOM_SECRET_KEY="your-secret-key" \
  -e KIOOM_MOCK="true" \
  -e KIOOM_MCP_SSE_TOKEN="my-secret-token" \
  --name kioom-mcp-container \
  kioom-mcp-sse
```

#### 3. GitHub Container Registry (GHCR) 이미지 사용
GitHub Actions를 통해 자동으로 빌드 및 배포되는 멀티 아키텍처(`amd64`, `arm64`) 공식 이미지를 직접 풀(pull)받아 실행할 수도 있습니다.
```bash
docker run -d \
  -p 8765:8765 \
  -e KIOOM_APP_KEY="your-app-key" \
  -e KIOOM_SECRET_KEY="your-secret-key" \
  -e KIOOM_MOCK="true" \
  -e KIOOM_MCP_SSE_TOKEN="my-secret-token" \
  --name kioom-mcp-container \
  ghcr.io/suapapa/go_kioom/kioom-mcp:latest
```

인증 토큰이 설정된 경우, 클라이언트는 다음과 같이 인증하여 접근해야 합니다:

1. **HTTP Header (추천)**:
   ```bash
   curl -I -H "Authorization: Bearer my-secret-token" http://127.0.0.1:8765/
   ```

2. **URL Query Parameter (EventSource 등 헤더 설정이 제한될 때)**:
   ```bash
   curl -I "http://127.0.0.1:8765/?token=my-secret-token"
   # 또는
   curl -I "http://127.0.0.1:8765/?auth=my-secret-token"
   ```

클라이언트는 이 엔드포인트에 대해 MCP SSE 규격(세션 `GET`, 메시지 `POST` 등)을 지원해야 합니다. 로컬이 아닌 주소로 노출할 때는 방화벽·TLS·인증을 반드시 고려하세요.

## 노출되는 도구(MCP tools)

이름은 MCP 식별자 규칙에 맞춰 밑줄(`_`)을 사용합니다.

| 도구 이름 | 설명 |
|-----------|------|
| `auth_issue` | 액세스 토큰 발급 |
| `auth_revoke` | 토큰 폐기 |
| `account_number` | 계좌번호 조회 |
| `account_deposit` | 예수금 등 (`kioom.DepositRequest` JSON) |
| `account_balance` | 잔고 조회 (`kioom.AccountBalanceRequest`) |
| `stock_basic` | 종목 기본 정보 (`stk_cd` 등) |
| `stock_rank` | 실시간 조회 순위 (`kioom.RealtimeItemRankRequest`) |
| `stock_minute_chart` | 분봉 (`kioom.StockMinuteChartRequest`) |
| `order_buy` / `order_sell` | 매수·매도 주문 |
| `order_modify` / `order_cancel` | 정정·취소 |

요청 필드 검증은 CLI와 공유하는 `internal/kioomvalidate`를 사용합니다.

## 주의 (거래·보안)

- 이 서버는 **실제 주문 API**에 연결될 수 있습니다. 운영·모의 환경과 키를 명확히 구분하고, 에이전트에게 부여할 권한을 제한하세요.
- SSE를 `0.0.0.0` 등으로 열면 외부에서 접근될 수 있습니다. 가능하면 `127.0.0.1`만 사용하거나 리버스 프록시 뒤에 두세요.

## 구현 위치

- 이 디렉터리: `main`, 플래그·환경 로딩·전송 선택
- [`internal/mcpkioom`](../../internal/mcpkioom): `mcp.Server`에 도구 등록

새 Kiwoom API를 MCP로 노출할 때는 **`internal/mcpkioom`만 수정**하면 stdio와 SSE가 함께 반영됩니다.
