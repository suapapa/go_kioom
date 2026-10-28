# go_kioom : Go 언어로 작성한 키움 증권 도구 모음

## MCP (`kioom-mcp`)

동일한 `Client`를 Model Context Protocol(MCP) 서버로 노출합니다. Cursor, Claude Desktop 같은 MCP 클라이언트는 stdio(기본)나 HTTP+SSE로 연결하면 됩니다. 인증·검증 규칙은 `kioom-cli`와 같습니다. 설치, 환경 변수, 노출 도구 목록은 [cmd/kioom-mcp/README.md](cmd/kioom-mcp/README.md)를 보세요.

## CLI (`kioom-cli`)

JSON 입출력과 스키마 조회용 CLI입니다. 설치, 옵션, 오류 코드, 출력 규약은 [cmd/kioom-cli/README.md](cmd/kioom-cli/README.md)를 참고하세요.

## Go package for Kiwoom REST API

[![Go Reference](https://pkg.go.dev/badge/github.com/suapapa/go_kioom.svg)](https://pkg.go.dev/github.com/suapapa/go_kioom)

> `go_kioom`은 키움증권 REST API 명세에 맞춘 비공식 Golang 래퍼(Wrapper) 패키지입니다.

### 사용 방법 (예제)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/suapapa/go_kioom"
)

func main() {
	appKey := "YOUR_APP_KEY"
	secretKey := "YOUR_SECRET_KEY"

	// 1. 클라이언트 생성 (모의투자는 WithMockDomain, 토큰 자동 발급은 WithAutoToken)
	client := kioom.NewClient(appKey, secretKey, kioom.WithAutoToken(), kioom.WithMockDomain())
	ctx := context.Background()

	// 2. 계좌번호 조회
	acctRes, err := client.GetAccountNumber(ctx)
	if err != nil {
		log.Printf("계좌번호 조회 실패: %v", err)
	} else {
		fmt.Printf("계좌번호: %s\n", acctRes.AcctNo)
	}

	// 3. 삼성전자 주식 기본정보 조회
	stockRes, err := client.GetStockBasicInfo(ctx, "005930")
	if err != nil {
		log.Printf("주식정보 조회 실패: %v", err)
	} else {
		fmt.Printf("주식명: %s, 현재가: %s\n", stockRes.StkNm, stockRes.CurPrc)
	}

	// 4. 접근토큰 폐기
	if _, err := client.RevokeToken(ctx); err != nil {
		log.Printf("토큰 폐기 실패: %v", err)
	} else {
		fmt.Println("토큰 폐기 완료")
	}
}
```

### 유량 제어 (Rate limits)

키움 REST API는 TR별로 초당 호출 수를 제한합니다. `go_kioom`은 제한이 공지된 TR을 클라이언트에서 자동으로 조절하므로, 루프나 여러 고루틴에서 호출해도 유량제어 오류를 피할 수 있습니다.

- `kt00018` (계좌평가잔고내역요청, `GetAccountBalance`): **초당 3회** — 2026-10-29 18:00 KST부터 초당 5회에서 3회로 강화
- 그 외 TR: 별도 제한 없음(서버 제한에 따름)

```go
// 기본값 사용 (kt00018은 초당 3회로 자동 조절)
client := kioom.NewClient(appKey, secretKey, kioom.WithAutoToken())

// 특정 TR의 제한을 직접 지정
client = kioom.NewClient(appKey, secretKey, kioom.WithRateLimit("ka10001", 5))

// 호출 측에서 이미 조절 중이라면 클라이언트 제한 해제
client = kioom.NewClient(appKey, secretKey, kioom.WithoutRateLimit())
```

클라이언트 하나를 공유하면 동시 호출에도 제한이 함께 적용됩니다. 대기 중 `ctx`가 취소되면 `context.Canceled`/`context.DeadlineExceeded`가 반환됩니다.

### 예제 (Examples)

더 많은 예제는 [cmd/examples/](./cmd/examples/)에서 볼 수 있습니다.

- `01_auth`: 토큰 발급 및 폐기
- `02_basic_info`: 주식 기본 정보 및 실시간 검색 순위 조회
- `03_monitor`: 계좌 현황 모니터링
- `04_trading`: 주식/신용 매수, 매도, 정정 및 취소 주문
- `05_minute_chart`: 주식 분봉 차트 조회

## 구현 현황

Kiwoom REST API 전체 구현 목록과 진행 상태는 [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md)에서 확인할 수 있습니다. API 필드·TR 코드는 공식 스펙 [`_ref/kiwoom-rest-api-spec.json`](_ref/kiwoom-rest-api-spec.json)을 기준으로 합니다.

### 확장 가이드

`kioom` 패키지는 새 API 엔드포인트를 쉽게 붙일 수 있게 `Client` 안에 공통 HTTP 요청 로직(`c.newRequest`, `c.do`)을 넣어 두었습니다. 공식 문서를 보고 Request, Response 구조체를 만든 뒤 메서드를 추가하면 됩니다.
